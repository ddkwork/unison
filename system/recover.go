package system

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"cogentcore.org/core/base/errors"
	"github.com/ddkwork/golibrary/mylog"
)

var HandleRecover = func(r any) {
	HandleRecoverBase(r)
	HandleRecoverPanic(r)
}

func HandleRecoverBase(r any) {
	if r == nil {
		return
	}
	stack := string(debug.Stack())
	log.Println("panic:", r)
	log.Println("")
	log.Println("----- START OF STACK TRACE: -----")
	log.Println(stack)
	log.Println("----- END OF STACK TRACE -----")

	dnm := filepath.Join(TheApp.AppDataDir(), "crash-logs")
	mylog.Check(os.MkdirAll(dnm, 0755))
	if errors.Log(err) != nil {
		return
	}
	cfnm := filepath.Join(dnm, "crash_"+time.Now().Format("2006-01-02_15-04-05"))
	mylog.Check(os.WriteFile(cfnm, []byte(CrashLogText(r, stack)), 0666))
	if errors.Log(err) != nil {
		return
	}
	cfnm = strings.ReplaceAll(cfnm, " ", `\ `)
	log.Println("SAVED CRASH LOG TO", cfnm)
}

func HandleRecoverPanic(r any) {
	if r == nil {
		return
	}
	if !TheApp.Platform().IsMobile() || TheApp.Platform() == Web {
		panic(r)
	}
}

func CrashLogText(r any, stack string) string {
	info := TheApp.SystemInfo()
	if info != "" {
		info += "\n"
	}
	return fmt.Sprintf("Platform: %v\nSystem platform: %v\nApp version: %s\nCore version: %s\nTime: %s\n%s\npanic: %v\n\n%s", TheApp.Platform(), TheApp.SystemPlatform(), AppVersion, CoreVersion, time.Now().Format(time.DateTime), info, r, stack)
}
