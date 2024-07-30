//go:build windows

package desktop

import (
	"os/exec"
	"os/user"
	"path/filepath"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/fileinfo/mimedata"
	"cogentcore.org/core/system"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/glfw"
)

func (a *App) Platform() system.Platforms {
	return system.Windows
}

func (a *App) OpenURL(url string) {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	errors.Log(cmd.Run())
}

func (a *App) DataDir() string {
	usr := mylog.Check2(user.Current())
	if errors.Log(err) != nil {
		return "/tmp"
	}
	return filepath.Join(usr.HomeDir, "AppData", "Roaming")
}

var TheClipboard = &Clipboard{}

type Clipboard struct{}

func (cl *Clipboard) IsEmpty() bool {
	str := glfw.GetClipboardString()
	return len(str) == 0
}

func (cl *Clipboard) Read(types []string) mimedata.Mimes {
	str := glfw.GetClipboardString()
	if len(str) == 0 {
		return nil
	}
	wantText := mimedata.IsText(types[0])
	if wantText {
		bstr := []byte(str)
		isMulti, mediaType, boundary, body := mimedata.IsMultipart(bstr)
		if isMulti {
			return mimedata.FromMultipart(body, boundary)
		} else {
			if mediaType != "" {
				return mimedata.NewMime(mediaType, bstr)
			} else {
				return mimedata.NewMime(types[0], bstr)
			}
		}
	} else {
	}
	return nil
}

func (cl *Clipboard) Write(data mimedata.Mimes) error {
	if len(data) == 0 {
		return nil
	}

	if len(data) > 1 {
		mpd := data.ToMultipart()
		glfw.SetClipboardString(string(mpd))
	} else {
		d := data[0]
		if mimedata.IsText(d.Type) {
			glfw.SetClipboardString(string(d.Data))
		}
	}
	return nil
}

func (cl *Clipboard) Clear() {
}
