// Package services
package services

import (
	"os/exec"
	"strings"

	"changeme/internal/models"
)

func Launch(item models.SearchResult) {
	go func() {
		var cmd *exec.Cmd
		if strings.HasSuffix(item.Path, ".desktop") {
			cmd = exec.Command("gtk-launch", item.Path)
		} else {
			cmd = exec.Command("xdg-open", item.Path)
		}

		_ = cmd.Start()
	}()
}

func RunCommand(command string, password string) string {
	var cmd *exec.Cmd
	if password != "" {
		cmd = exec.Command("sh", "-c", "echo '"+password+"' | sudo -S "+command+" 2>&1")
	} else {
		cmd = exec.Command("sh", "-c", command+" 2>&1")
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out) + "\nError: " + err.Error()
	}
	return string(out)
}
