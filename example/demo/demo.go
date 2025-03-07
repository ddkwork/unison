// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package demo

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/ddkwork/unison/enums/align"
	"github.com/ddkwork/unison/enums/behavior"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/toolbox/errs"
	"github.com/ddkwork/unison"
)

var windowCounter int

// NewDemoWindow creates and displays our demo window.
func NewDemoWindow(where unison.Point) (*unison.Window, error) {
	// Create the window
	windowCounter++
	wnd := mylog.Check2(unison.NewWindow(fmt.Sprintf("Demo #%d", windowCounter)))

	// Install our menus
	installDefaultMenus(wnd)

	// Put some empty space around the edges of our window and apply a single column layout.
	content := wnd.Content()
	content.SetBorder(unison.NewEmptyBorder(unison.NewUniformInsets(10)))
	content.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: 10,
	})

	// Create a wrappable row of buttons
	panel := createButtonsPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	// Create a wrappable row of buttons that bring up dialogs
	panel = createDialogButtonsPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	// Create a wrappable row of image buttons
	panel = createImageButtonsPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	addSeparator(content)

	// Create a column of checkboxes
	panel = createCheckBoxPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	addSeparator(content)

	// Create a column of radio buttons and a progress bar they control
	panel = createRadioButtonsAndProgressBarsPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	addSeparator(content)

	// Create a column of popup menus
	panel = createPopupMenusPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	addSeparator(content)

	// Create some fields and a list, side-by-side
	panel = createFieldsAndListPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	addSeparator(content)

	// Create an image panel, but don't add it yet
	imgPanel := createImagePanel()

	// Create some color wells and pass it our image panel
	panel = createWellsPanel(imgPanel)
	panel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		VAlign: align.Middle,
		HGrab:  true,
	})
	content.AddChild(panel)

	// Create a scroll panel and place the image panel inside it
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(imgPanel, behavior.Unmodified, behavior.Unmodified)
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	content.AddChild(scrollArea)

	// Marshal our window to fit its content, then set its location on the display and make it visible.
	wnd.Pack()
	rect := wnd.FrameRect()
	rect.Point = where
	wnd.SetFrameRect(rect)
	wnd.ToFront()

	return wnd, nil
}

func createButtonsPanel() *unison.Panel {
	// Create a panel to place some buttons into.
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})

	// Add some buttons
	for i, title := range []string{"First", "Second", "Third", "Fourth", "Fifth", "Sixth", "Seventh", "Eighth", "Ninth"} {
		btn := createButton(title, panel)
		if i == 2 {
			btn.SetEnabled(false)
		}
	}
	return panel
}

func createDialogButtonsPanel() *unison.Panel {
	// Create a panel to place some buttons into.
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})

	btn := createButton("Show Question Dialog", panel)
	btn.ClickCallback = func() {
		unison.QuestionDialog("Sample Question", "Text for a question goes here")
	}
	btn = createButton("Show Warning Dialog", panel)
	btn.ClickCallback = func() {
		unison.WarningDialogWithMessage("Sample Warning", "Text for a warning goes here")
	}
	btn = createButton("Show Error Dialog", panel)
	btn.ClickCallback = func() {
		unison.ErrorDialogWithError("Sample Error", errs.New("A stack trace will be emitted to the log"))
	}

	return panel
}

func createButton(title string, panel *unison.Panel) *unison.Button {
	btn := unison.NewButton()
	btn.Text = title
	btn.ClickCallback = func() { slog.Info(title) }
	btn.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("Tooltip for: %s", title))
	btn.SetLayoutData(align.Middle)
	panel.AddChild(btn)
	return btn
}

func createImageButtonsPanel() *unison.Panel {
	// Create a panel to place some buttons into.
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})

	// Load our home image, and if successful (we should be!), add two buttons with it, one enabled and one not.
	homeImg := mylog.Check2(HomeImage())

	// Load our logo image, and if successful (we should be!), add two buttons with it, one enabled and one not.
	var logoImg *unison.Image
	logoImg = mylog.Check2(ClassicAppleLogoImage())
	createImageButton(logoImg, "logo_enabled", panel)
	createImageButton(logoImg, "logo_disabled", panel).SetEnabled(false)

	if homeImg != nil && logoImg != nil {
		// Add spacer
		spacer := &unison.Panel{}
		spacer.Self = spacer
		spacer.SetSizer(func(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
			minSize.Width = 40
			prefSize.Width = 40
			maxSize.Width = 40
			return
		})
		panel.AddChild(spacer)

		// Add some sticky buttons in a group with our images
		group := unison.NewGroup()
		first := createImageButton(homeImg, "home_toggle", panel)
		first.Sticky = true
		group.Add(first.AsGroupPanel())
		second := createImageButton(logoImg, "logo_toggle", panel)
		second.Sticky = true
		group.Add(second.AsGroupPanel())
		group.Select(first.AsGroupPanel())
	}

	return panel
}

func createImageButton(img *unison.Image, actionText string, panel *unison.Panel) *unison.Button {
	btn := unison.NewButton()
	btn.Drawable = img
	btn.ClickCallback = func() { slog.Info(actionText) }
	btn.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("Tooltip for: %s", actionText))
	btn.SetLayoutData(align.Middle)
	panel.AddChild(btn)
	return btn
}

func addSeparator(parent *unison.Panel) {
	sep := unison.NewSeparator()
	sep.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
	})
	parent.AddChild(sep)
}

func createCheckBoxPanel() *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	createCheckBox("Initially Off", unison.OffCheckState, panel)
	createCheckBox("Initially On", unison.OnCheckState, panel)
	createCheckBox("Initially Mixed", unison.MixedCheckState, panel)
	createCheckBox("Disabled", unison.OffCheckState, panel).SetEnabled(false)
	createCheckBox("Disabled w/Check", unison.OnCheckState, panel).SetEnabled(false)
	createCheckBox("Disabled w/Mixed", unison.MixedCheckState, panel).SetEnabled(false)
	return panel
}

func createCheckBox(title string, initialState unison.CheckState, panel *unison.Panel) *unison.CheckBox {
	check := unison.NewCheckBox()
	check.Text = title
	check.State = initialState
	check.ClickCallback = func() { slog.Info("checkbox clicked", "title", title) }
	check.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("This is the tooltip for '%s'", title))
	panel.AddChild(check)
	return check
}

func createRadioButtonsAndProgressBarsPanel() *unison.Panel {
	// Create a wrapper to put them side-by-side
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:      2,
		HSpacing:     10,
		VSpacing:     unison.StdVSpacing,
		VAlign:       align.Middle,
		EqualColumns: true,
	})

	// Create the progress bar, but don't add it yet
	progress := unison.NewProgressBar(100)
	progress.SetCurrent(25)
	progress.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})

	// Create the radio buttons that will control the progress bar
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	group := unison.NewGroup()
	first := createRadioButton("25%", panel, group, progress, 25, 100)
	createRadioButton("50%", panel, group, progress, 50, 100)
	createRadioButton("75%", panel, group, progress, 75, 100).SetEnabled(false)
	createRadioButton("100%", panel, group, progress, 100, 100)
	createRadioButton("Indeterminate", panel, group, progress, 0, 0)
	group.Select(first.AsGroupPanel())

	// Add the radio buttons to the left
	wrapper.AddChild(panel)

	// Add the progress bar to the right
	wrapper.AddChild(progress)

	return wrapper
}

func createRadioButton(title string, panel *unison.Panel, group *unison.Group, progressBar *unison.ProgressBar, current, maximum float32) *unison.RadioButton {
	rb := unison.NewRadioButton()
	rb.Text = title
	rb.ClickCallback = func() {
		progressBar.SetMaximum(maximum)
		progressBar.SetCurrent(current)
		slog.Info("radio button clicked", "title", title)
	}
	rb.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("This is the tooltip for %s", title))
	panel.AddChild(rb)
	group.Add(rb.AsGroupPanel())
	return rb
}

func createPopupMenusPanel() *unison.Panel {
	panel := unison.NewPanel()
	createPopupMenu(panel, 1, "Alphabet Tooltip", "Alpha", "Beta", "Charlie", "", "Delta", "Echo", "Foxtrot")
	createPopupMenu(panel, 2, "Color Tooltip", "Red", "Blue", "Green").SetEnabled(false)
	panel.SetLayout(&unison.FlexLayout{
		Columns:  len(panel.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	return panel
}

func createPopupMenu(panel *unison.Panel, selection int, tooltip string, titles ...string) *unison.PopupMenu[string] {
	p := unison.NewPopupMenu[string]()
	p.Tooltip = unison.NewTooltipWithText(tooltip)
	for _, title := range titles {
		if title == "" {
			p.AddSeparator()
		} else {
			p.AddItem(title)
		}
	}
	p.SelectIndex(selection)
	p.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
		if title, ok := popup.Selected(); ok {
			slog.Info("item selected from PopupMenu", "title", title, "popup", tooltip)
		}
	}
	panel.AddChild(p)
	return p
}

func createFieldsAndListPanel() *unison.Panel {
	// Create a wrapper to put them side-by-side
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:      2,
		HSpacing:     10,
		VSpacing:     unison.StdVSpacing,
		EqualColumns: true,
	})

	// Add the text fields to the left side
	textFieldsPanel := createTextFieldsPanel()
	textFieldsPanel.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	wrapper.AddChild(textFieldsPanel)

	// Add the list to the right side
	wrapper.AddChild(createListPanel())

	return wrapper
}

func createTextFieldsPanel() *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	createTextField("Field 1:", "First Text Field", panel)
	createTextField("Field 2:", "Second Text Field (disabled)", panel).SetEnabled(false)
	field := createTextField("Longer Label:", "", panel)
	field.Watermark = "Password Field"
	field.ObscurementRune = '●'
	field = createTextField("Field 4:", "", panel)
	field.HAlign = align.End
	field.Watermark = "Enter only numbers"
	field.ValidateCallback = func() bool {
		for _, r := range field.Text() {
			if !unicode.IsDigit(r) {
				return false
			}
		}
		return true
	}
	createMultiLineTextField("Field 5:", "One\nTwo\nThree", panel)
	return panel
}

func createTextField(labelText, fieldText string, panel *unison.Panel) *unison.Field {
	lbl := unison.NewLabel()
	lbl.Text = labelText
	lbl.HAlign = align.End
	lbl.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.End,
		VAlign: align.Middle,
	})
	panel.AddChild(lbl)
	field := unison.NewField()
	field.SetText(fieldText)
	field.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("This is the tooltip for %v", field))
	panel.AddChild(field)
	return field
}

func createMultiLineTextField(labelText, fieldText string, panel *unison.Panel) *unison.Field {
	lbl := unison.NewLabel()
	lbl.Text = labelText
	lbl.HAlign = align.End
	lbl.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.End,
		VAlign: align.Middle,
	})
	panel.AddChild(lbl)
	field := unison.NewMultiLineField()
	field.SetText(fieldText)
	field.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Middle,
		HGrab:  true,
	})
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("This is the tooltip for %v", field))
	panel.AddChild(field)
	return field
}

func createListPanel() *unison.Panel {
	lst := unison.NewList[string]()
	lst.Append(
		"One",
		"Two",
		"Three with some long text to make it interesting",
		"Four",
		"Five",
	)
	lst.NewSelectionCallback = func() {
		var buffer strings.Builder
		buffer.WriteString("Selection changed in the list. Now:")
		index := -1
		first := true
		for {
			index = lst.Selection.NextSet(index + 1)
			if index == -1 {
				break
			}
			if first {
				first = false
			} else {
				buffer.WriteString(",")
			}
			fmt.Fprintf(&buffer, " %d", index)
		}
		slog.Info(buffer.String())
	}
	lst.DoubleClickCallback = func() {
		slog.Info("Double-clicked on the list")
	}
	_, prefSize, _ := lst.Sizes(unison.Size{})
	lst.SetFrameRect(unison.Rect{Size: prefSize})
	scroller := unison.NewScrollPanel()
	scroller.SetBorder(unison.NewLineBorder(unison.ControlEdgeColor, 0, unison.NewUniformInsets(1), false))
	scroller.SetContent(lst, behavior.Fill, behavior.Fill)
	scroller.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	return scroller.AsPanel()
}

func createImagePanel() *unison.Label {
	// Create the label and make it focusable
	imgPanel := unison.NewLabel()
	imgPanel.SetFocusable(true)

	// Prepare a cursor for when the mouse is over the image
	cursor := unison.MoveCursor()
	logoImg := mylog.Check2(ClassicAppleLogoImage())
	size := logoImg.LogicalSize()
	cursor = unison.NewCursor(logoImg, unison.Point{
		X: size.Width / 2,
		Y: size.Height / 2,
	})
	imgPanel.UpdateCursorCallback = func(where unison.Point) *unison.Cursor { return cursor }

	// Add a tooltip that shows the current mouse coordinates
	imgPanel.UpdateTooltipCallback = func(where unison.Point, suggestedAvoidInRoot unison.Rect) unison.Rect {
		imgPanel.Tooltip = unison.NewTooltipWithText(where.String())
		suggestedAvoidInRoot.X = where.X - 16
		suggestedAvoidInRoot.Y = where.Y - 16
		suggestedAvoidInRoot.Width = 32
		suggestedAvoidInRoot.Height = 32
		return imgPanel.RectToRoot(suggestedAvoidInRoot)
	}

	// Set the initial image
	mylog.Check2(MountainsImage())

	// Set the set of the widget to match its preferred size
	_, prefSize, _ := imgPanel.Sizes(unison.Size{})
	imgPanel.SetFrameRect(unison.Rect{Size: prefSize})

	return imgPanel
}

func createWellsPanel(imgPanel *unison.Label) *unison.Panel {
	// Create the panel that's going to hold the wells
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: 5,
		VSpacing: 5,
	})

	// Add a well
	well1 := unison.NewWell()
	well1.SetInk(unison.Yellow)
	panel.AddChild(well1)

	// When this well is changed, if the user sets an image, we'll change the image panel to match it
	well1.InkChangedCallback = func() {
		ink := well1.Ink()
		if pattern, ok := ink.(*unison.Pattern); ok {
			imgPanel.Drawable = pattern.Image
			_, pSize, _ := imgPanel.Sizes(unison.Size{})
			imgPanel.SetFrameRect(unison.Rect{Size: pSize})
			imgPanel.MarkForRedraw()
		}
	}

	// Add another, disabled, well
	well2 := unison.NewWell()
	well2.SetInk(unison.Orange)
	well2.SetEnabled(false)

	panel.AddChild(well2)
	return panel
}
