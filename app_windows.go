// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package unison

import (
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ddkwork/unison/enums/thememode"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/unison/internal/w32"
	"golang.org/x/sys/windows/registry"
)

var appUsesLightThemeValue = uint32(1)

func platformEarlyInit() {
	AttachConsole()
}

func platformLateInit() {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	k := mylog.Check2(registry.OpenKey(registry.CURRENT_USER, keyPath, syscall.KEY_NOTIFY|registry.QUERY_VALUE))

	mylog.Check(updateTheme(k, true))
	go func() {
		for {
			w32.RegNotifyChangeKeyValue(k, false, w32.RegNotifyChangeName|w32.RegNotifyChangeLastSet, 0, false)
			mylog.Check(updateTheme(k, false))
		}
	}()
}

func platformFinishedStartup() {
}

func platformBeep() {
	w32.MessageBeep(w32.MBDefault)
}

func platformIsDarkModeTrackingPossible() bool {
	return true
}

func platformIsDarkModeEnabled() bool {
	return atomic.LoadUint32(&appUsesLightThemeValue) == 0
}

func platformDoubleClickInterval() time.Duration {
	return w32.GetDoubleClickTime()
}

func updateTheme(k registry.Key, sync bool) error {
	val, _ := mylog.Check3(k.GetIntegerValue("AppsUseLightTheme"))

	var swapped bool
	if val == 0 {
		swapped = atomic.CompareAndSwapUint32(&appUsesLightThemeValue, 1, 0)
	} else {
		swapped = atomic.CompareAndSwapUint32(&appUsesLightThemeValue, 0, 1)
	}
	if swapped && currentColorMode == thememode.Auto {
		if sync {
			ThemeChanged()
		} else {
			InvokeTask(ThemeChanged)
		}
	}
	return nil
}
