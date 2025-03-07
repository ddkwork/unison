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
	"sync"
	"time"

	"github.com/ddkwork/toolbox"
	"github.com/ddkwork/toolbox/errs"
)

var (
	taskQueueLock sync.Mutex
	taskQueue     []func()
)

// InvokeTask calls a function on the UI thread. The function is put into the system event queue and will be run at the
// next opportunity.
func InvokeTask(f func()) {
	taskQueueLock.Lock()
	taskQueue = append(taskQueue, f)
	taskQueueLock.Unlock()
	postEmptyEvent()
}

// InvokeTaskAfter schedules a function to be run on the UI thread after waiting for the specified duration.
func InvokeTaskAfter(f func(), after time.Duration) {
	time.AfterFunc(after, func() { InvokeTask(f) })
}

func processNextTask(recoveryHandler errs.RecoveryHandler) {
	var f func()
	needsPost := false
	taskQueueLock.Lock()
	if len(taskQueue) > 0 {
		f = taskQueue[0]
		copy(taskQueue, taskQueue[1:])
		taskQueue = taskQueue[:len(taskQueue)-1]
		needsPost = len(taskQueue) > 0
	}
	taskQueueLock.Unlock()
	if f != nil {
		toolbox.CallWithHandler(f, recoveryHandler)
		if needsPost {
			postEmptyEvent()
		}
	}
}
