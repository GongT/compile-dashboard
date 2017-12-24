package gui

import (
	"github.com/jroimartin/gocui"
	"errors"
	"github.com/gongt/compile-dashboard/lib"
	"github.com/nsf/termbox-go"
)

type Control struct {
	g       *gocui.Gui
	tabs    *tabManager
	state   *tabPanel
	actions map[string]initAction
}

func NewControl() (ui *Control, err error) {
	defer func() {
		if err != nil {
			termboxClose()
		}
	}()

	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		lib.MainLogger.Println("GUI init failed", err)
		return
	}
	gui.Cursor = true
	gui.Mouse = true

	tab := newTabManager(gui)
	state := newTabPanel(gui, tab, func(i int) {
		tab.switchTo(i)
	})

	ui = &Control{
		gui,
		tab,
		state,
		make(map[string]initAction),
	}

	gui.SetManagerFunc(func(_ *gocui.Gui) error {
		return ui.render(gui)
	})

	err = state.initKeys()
	if err != nil {
		return
	}

	err = ui.initKeys()
	if err != nil {
		return
	}

	return
}

func termboxClose() {
	if !termbox.IsInit {
		lib.MainLogger.Println("called terminate screen (but screen already stopped).")
		return
	}

	lib.MainLogger.Println("called terminate screen.")
	termbox.Close()
	lib.MainLogger.Println("screen terminated.")
}

func (ui *Control) Close() {
	lib.MainLogger.Println("called UI.Close().")
	termboxClose()
}

func (ui *Control) initKeys() error {
	gui := ui.g
	if err := gui.SetKeybinding(viewNameAll, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *Control) render(g *gocui.Gui) error {
	event := newRenderEvent(g)
	var err error

	err = subLayout(g, viewNameMessage, Geo{event.splitCenter, event.splitSideMiddle, event.maxX - 1, event.maxY - 2}, func(view *gocui.View) {
		view.Wrap = true
		view.Title = "Message"
	})
	if err != nil {
		return err
	}

	err = subLayout(g, viewNameBottom, Geo{0, event.maxY - 2, event.maxX, event.maxY}, func(view *gocui.View) {
		view.Frame = false
	})
	if err != nil {
		return err
	}

	err = ui.state.render(event)
	if err != nil {
		return err
	}

	return ui.tabs.render(event)
}

func (ui *Control) EventLoop() error {
	lib.MainLogger.Println("tab count=", ui.tabs.length)

	if ui.tabs.length == 0 {
		ui.g.Close()
		lib.MainLogger.Println("raising error: no tabs")
		return errors.New("no Tabs, can not start event loop")
	}
	if err := ui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (ui *Control) AddTab(name string, action initAction) {
	lib.MainLogger.Println("add tab", name)
	ui.tabs.add(name, action)
	ui.actions[name] = action
}

func (ui *Control) MarkTabError(name string, hasError bool) error {
	index, err := ui.tabs.findIndex(name)
	if err != nil {
		return err
	}
	ui.state.setErrorStatus(index, hasError)
	return nil
}

func (ui *Control) Update(view *gocui.View) error {
	i, e := ui.tabs.findIndex(view.Name())
	if i == -1 {
		return e
	}
	if i == ui.tabs.current {
		ui.g.Update(emptyRenderEvent)
		lib.MainLogger.Println("update view:", view.Name())
	}

	return nil
}
