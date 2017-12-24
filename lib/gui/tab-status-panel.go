package gui

import (
	"github.com/jroimartin/gocui"
	"github.com/gongt/compile-dashboard/lib"
)

type currentTab int

type tabPanel struct {
	gui      *gocui.Gui
	tab      *tabManager
	hasErr   [maxTabCount]bool
	activate currentTab
	callback func(int)
}

func newTabPanel(gui *gocui.Gui, tabControl *tabManager, callback func(int)) *tabPanel {
	ret := new(tabPanel)

	ret.gui = gui
	ret.tab = tabControl
	ret.hasErr = [maxTabCount]bool{}
	ret.activate = 0
	ret.callback = callback

	return ret
}

func (tp *tabPanel) setErrorStatus(tabIndex int, has bool) {
	tp.hasErr[tabIndex] = has
	tp.gui.Update(func(gui *gocui.Gui) error {
		view, err := gui.View(viewNameStatus)
		if err != nil {
			return err
		}
		return tp.update(view)
	})
}

func (tp *tabPanel) initKeys() error {
	return tp.gui.SetKeybinding(viewNameStatus, gocui.MouseLeft, gocui.ModNone, func(_ *gocui.Gui, view *gocui.View) error {
		_, y := view.Cursor()

		if tp.activate == currentTab(y) {
			return nil
		}
		lib.MainLogger.Printf("switch to tab %d/%d\n", y, tp.tab.length)
		if y >= tp.tab.length {
			return nil
		}
		lib.MainLogger.Printf("switch tab %d -> %d\n", tp.activate, y)

		tp.activate = currentTab(y)
		tp.callback(y)
		return tp.update(view)
	})
}

func (tp *tabPanel) render(e *renderEvent) error {
	return subLayout(e.g, viewNameStatus, Geo{e.splitCenter, 0, e.maxX - 1, e.splitSideMiddle - 1}, func(view *gocui.View) {
		view.Wrap = false
		view.Title = "Status"
		tp.update(view)
	})
}

func (tp *tabPanel) update(view *gocui.View) error {
	view.Clear()
	for i, v := range tp.tab.tabNames {
		if v == nil {
			break
		}
		var color string
		if currentTab(i) == tp.activate {
			if tp.hasErr[i] {
				color = "9"
			} else {
				color = "10"
			}
			view.Write([]byte("\x1B[38;5;0m\x1B[48;5;" + color + "m"))
		} else if tp.hasErr[i] {
			view.Write([]byte("\x1B[0m\x1B[38;5;9m"))
		} else {
			view.Write([]byte("\x1B[0m"))
		}

		view.Write([]byte(*v))
		view.Write([]byte("                       \n"))
	}

	return nil
}
