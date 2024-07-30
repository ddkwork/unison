//go:build darwin

package desktop

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/fsnotify/fsnotify"
)

const plistPath = `/Library/Preferences/.GlobalPreferences.plist`

var plist = filepath.Join(os.Getenv("HOME"), plistPath)

func (a *App) IsDark() bool {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	if mylog.Check(cmd.Run()); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false
		} else {
			slog.Error("unexpected error when running command to get system color theme: " + err.Error())
			return false
		}
	}
	return true
}

func (a *App) IsDarkMonitor() {
	watcher := mylog.Check2(fsnotify.NewWatcher())
	mylog.Check(watcher.Add(plist))

	defer watcher.Close()
	wasDark := a.IsDark()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				isDark := a.IsDark()
				if isDark != wasDark {
					a.Dark = isDark
					wasDark = isDark
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("system color theme watcher error: " + err.Error())
		}
	}
}
