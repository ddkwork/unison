// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package main

import (
	_ "embed"
	"os"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/toolbox/cmdline"
	"github.com/ddkwork/unison"
	"github.com/ddkwork/unison/example/demo"
)

func main() {
	cmdline.AppName = "Example"
	cmdline.AppCmdName = "example"
	cmdline.CopyrightStartYear = "2021"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.unison.example"

	unison.AttachConsole()

	cl := cmdline.New(true)
	cl.Parse(os.Args[1:])
	unison.Start(unison.StartupFinishedCallback(func() {
		mylog.Check2(demo.NewDemoTableWindow(unison.PrimaryDisplay().Usable.Point))
		// mylog.Check2(demo.NewDemoWindow(unison.PrimaryDisplay().Usable.Point))
	})) // Never returns
}
