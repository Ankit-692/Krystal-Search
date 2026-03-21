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
