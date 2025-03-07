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
	"github.com/ddkwork/toolbox/i18n"
	"github.com/ddkwork/unison/enums/align"
)

// DialogButtonInfo holds information for constructing the dialog button panel.
type DialogButtonInfo struct {
	Title        string
	ResponseCode int
	KeyCodes     []KeyCode
}

// NewButton creates a new button for the dialog.
func (bi *DialogButtonInfo) NewButton(d *Dialog) *Button {
	b := NewButton()
	b.Text = bi.Title
	b.ClickCallback = func() { d.StopModal(bi.ResponseCode) }
	b.SetLayoutData(&FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
	})
	return b
}

// NewCancelButtonInfo creates a standard cancel button.
func NewCancelButtonInfo() *DialogButtonInfo {
	return &DialogButtonInfo{
		Title:        i18n.Text("Cancel"),
		ResponseCode: ModalResponseCancel,
		KeyCodes:     []KeyCode{KeyEscape},
	}
}

// NewNoButtonInfo creates a standard No button.
func NewNoButtonInfo() *DialogButtonInfo {
	return &DialogButtonInfo{
		Title:        i18n.Text("No"),
		ResponseCode: ModalResponseDiscard,
		KeyCodes:     []KeyCode{KeyN},
	}
}

// NewYesButtonInfo creates a standard Yes button.
func NewYesButtonInfo() *DialogButtonInfo {
	return &DialogButtonInfo{
		Title:        i18n.Text("Yes"),
		ResponseCode: ModalResponseOK,
		KeyCodes:     []KeyCode{KeyReturn, KeyNumPadEnter, KeyY},
	}
}

// NewOKButtonInfo creates a standard OK button.
func NewOKButtonInfo() *DialogButtonInfo {
	return NewOKButtonInfoWithTitle(i18n.Text("OK"))
}

// NewOKButtonInfoWithTitle creates a standard OK button with a specific title.
func NewOKButtonInfoWithTitle(title string) *DialogButtonInfo {
	return &DialogButtonInfo{
		Title:        title,
		ResponseCode: ModalResponseOK,
		KeyCodes:     []KeyCode{KeyReturn, KeyNumPadEnter},
	}
}
