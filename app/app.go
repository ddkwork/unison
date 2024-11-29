package app

import (
	"log"
	"log/slog"

	"github.com/ddkwork/unison/enums/thememode"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/toolbox/log/tracelog"
	"github.com/ddkwork/unison"
)

/*
func (d *colorSettingsDockable) addToStartToolbar(toolbar *unison.Panel) {
	label := unison.NewLabel()
	label.SetTitle(i18n.Text("Color Mode"))
	toolbar.AddChild(label)
	p := unison.NewPopupMenu[thememode.EnumTypes]()
	for _, mode := range thememode.All {
		p.AddItem(mode)
	}
	p.Select(gurps.GlobalSettings().ThemeMode)
	p.SelectionChangedCallback = func(popup *unison.PopupMenu[thememode.EnumTypes]) {
		if mode, ok := popup.Selected(); ok {
			gurps.GlobalSettings().ThemeMode = mode
			unison.SetThemeMode(mode)
		}
	}
	toolbar.AddChild(p)
}

*/

func Run(title string, layoutCallback func(w *unison.Window)) {
	run(title, nil, layoutCallback)
}

func RunWithIco(title string, ico []byte, layoutCallback func(w *unison.Window)) {
	run(title, ico, layoutCallback)
}

func run(title string, ico []byte, layoutCallback func(w *unison.Window)) {
	mylog.Call(func() {
		unison.Start(unison.StartupFinishedCallback(func() {
			unison.SetThemeMode(thememode.Dark)
			w := mylog.Check2(unison.NewWindow(title))
			if ico != nil {
				b := mylog.Check2(unison.NewImageFromBytes(ico, 0.5))
				w.SetTitleIcons([]*unison.Image{b})
			}
			// installDefaultMenus(w)
			layoutCallback(w)
			w.Pack()

			rect := w.FrameRect()
			rect.Point = unison.PrimaryDisplay().Usable.Point
			if rect.Width < 200 {
				rect.Width = 800
			}
			if rect.Height < 10 {
				rect.Height = 600
			}
			rect = unison.BestDisplayForRect(rect).FitRectOnto(rect)
			// rect.Point = rect.Center()
			w.SetFrameRect(rect)

			w.ToFront()
		}))
	})
}

/*
42c26e1249b14c404a071392a2a6a51525d7b4d8 ok

d845b56964c1c227fde3bd819d5656818d43645a bug
*/

func init() {
	slog.SetDefault(slog.New(tracelog.New(log.Default().Writer(), slog.LevelInfo)))
	//errs.RuntimePrefixesToFilter = append(errs.RuntimePrefixesToFilter,
	//	"github.com/ddkwork/toolbox.callWithHandler",
	//	"github.com/ddkwork/toolbox.call",
	//)
}
