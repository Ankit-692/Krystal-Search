// Package services
package services

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatcherService struct {
	watcher *fsnotify.Watcher
	stopCh  chan struct{}
	once    sync.Once

	OnAppsChanged func()
	OnIconChanged func(newTheme string)
}

var appDirs = []string{
	"/usr/share/applications",
	filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
}

var gtkSettingsFiles = []string{
	filepath.Join(os.Getenv("HOME"), ".config", "gtk-3.0", "settings.ini"),
	filepath.Join(os.Getenv("HOME"), ".config", "gtk-4.0", "settings.ini"),
}

func NewWatcherService() (*WatcherService, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &WatcherService{
		watcher: w,
		stopCh:  make(chan struct{}),
	}, nil
}

func (ws *WatcherService) Start() {
	ws.watchAppDirs()
	ws.watchGTKSettingsDirs()

	go ws.listenFSEvents()
	go ws.pollGSettings()
}

func (ws *WatcherService) Stop() {
	ws.once.Do(func() {
		close(ws.stopCh)
		ws.watcher.Close()
	})
}

func (ws *WatcherService) watchAppDirs() {
	for _, dir := range appDirs {
		if _, err := os.Stat(dir); err != nil {
			log.Printf("[watcher]Not exists app dir %s: %v", dir, err)
			continue
		}
		if err := ws.watcher.Add(dir); err != nil {
			log.Printf("[watcher] cannot watch app dir %s: %v", dir, err)
		} else {
			log.Printf("[watcher] watching app dir: %s", dir)
		}
	}
}

func (ws *WatcherService) watchGTKSettingsDirs() {
	seen := map[string]bool{}
	for _, f := range gtkSettingsFiles {
		dir := filepath.Dir(f)
		if seen[dir] {
			continue
		}
		seen[dir] = true
		if _, err := os.Stat(dir); err != nil {
			log.Printf("[watcher] Dont Exists gtk dir %s: %v", dir, err)
			continue
		}
		if err := ws.watcher.Add(dir); err != nil {
			log.Printf("[watcher] cannot watch gtk dir %s: %v", dir, err)
		}
	}
}

func (ws *WatcherService) listenFSEvents() {
	appDebounce := newDebouncer(300 * time.Millisecond)

	// Snapshot of icon theme per GTK file so we only fire on real changes
	iconSnap := make(map[string]string)
	for _, f := range gtkSettingsFiles {
		iconSnap[f] = readIconThemeFromGTKFile(f)
	}
	iconDebounce := newDebouncer(400 * time.Millisecond)

	for {
		select {
		case <-ws.stopCh:
			return

		case event, ok := <-ws.watcher.Events:
			if !ok {
				return
			}

			switch {
			case isInAppDir(event.Name) && isDesktopFile(event.Name):
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) ||
					event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					log.Printf("[watcher] desktop file event: %s (%s)", event.Name, event.Op)
					appDebounce(func() {
						if ws.OnAppsChanged != nil {
							ws.OnAppsChanged()
						}
					})
				}

			case isGTKSettingsFile(event.Name):
				iconDebounce(func() {
					for _, f := range gtkSettingsFiles {
						theme := readIconThemeFromGTKFile(f)
						if theme != "" && theme != iconSnap[f] {
							log.Printf("[watcher] icon theme via GTK file: %s → %s", iconSnap[f], theme)
							iconSnap[f] = theme
							if ws.OnIconChanged != nil {
								ws.OnIconChanged(theme)
							}
						}
					}
				})
			}

		case err, ok := <-ws.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[watcher] fs error: %v", err)
		}
	}
}

func (ws *WatcherService) pollGSettings() {
	current := queryGSettingsTheme()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ws.stopCh:
			return
		case <-ticker.C:
			theme := queryGSettingsTheme()
			if theme != "" && theme != current {
				log.Printf("[watcher] icon theme via gsettings: %s → %s", current, theme)
				current = theme
				if ws.OnIconChanged != nil {
					ws.OnIconChanged(theme)
				}
			}
		}
	}
}

func queryGSettingsTheme() string {
	out, err := exec.Command("gsettings", "get",
		"org.gnome.desktop.interface", "icon-theme").Output()
	if err != nil {
		return "" // gsettings not available (KDE-only system, etc.)
	}
	// Output looks like: 'Papirus-Dark'\n
	return strings.Trim(strings.TrimSpace(string(out)), "'")
}

func isDesktopFile(path string) bool {
	return filepath.Ext(path) == ".desktop"
}

func isInAppDir(path string) bool {
	for _, dir := range appDirs {
		if strings.HasPrefix(path, dir) {
			return true
		}
	}
	return false
}

func isGTKSettingsFile(path string) bool {
	for _, f := range gtkSettingsFiles {
		if path == f {
			return true
		}
	}
	return false
}

func readIconThemeFromGTKFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "gtk-icon-theme-name") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func newDebouncer(d time.Duration) func(fn func()) {
	var (
		mu    sync.Mutex
		timer *time.Timer
	)
	return func(fn func()) {
		mu.Lock()
		defer mu.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(d, fn)
	}
}
