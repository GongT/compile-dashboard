package gui

import (
	"github.com/jroimartin/gocui"
	"github.com/gongt/compile-dashboard/lib"
	"github.com/nsf/termbox-go"
	"fmt"
	"sync"
)

type Control struct {
	gui       *gocui.Gui
	tabs      *tabManager
	state     *tabPanel
	waiting   chan bool
	Error     chan error
	writeLock sync.Mutex
	inited    bool
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
		make(chan bool),
		make(chan error),
		sync.Mutex{},
		false,
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

func (ui *Control) Inspect() string {
	return "the main screen controler"
}
func (ui *Control) Close() error {
	lib.MainLogger.Println("called UI.Close().")
	termboxClose()
	return nil
}

func (ui *Control) initKeys() error {
	gui := ui.gui
	if err := gui.SetKeybinding(viewNameAll, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *Control) WaitInit() {
	lib.MainLogger.Println("waiting screen to init: ")
	_, ok := <-ui.waiting
	lib.MainLogger.Println("this is first wait: ", ok)
}

func (ui *Control) render(g *gocui.Gui) error {
	event := newRenderEvent(ui)
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

	err = ui.tabs.render(event)
	if err != nil {
		return err
	}

	if !ui.inited {
		ui.inited = true
		lib.MainLogger.Println("first render, init complete!")
		ui.waiting <- true
		close(ui.waiting)
	}

	return nil
}

func (ui *Control) Message(msg string) {
	view, _ := ui.gui.View(viewNameMessage)
	fmt.Fprintln(view, msg)
}

func (ui *Control) StartEventLoop() {
	go func() {
		lib.MainLogger.Println("main event loop started")
		err := ui.gui.MainLoop()
		lib.MainLogger.Println("main event loop terminated")
		if err != nil {
			ui.Error <- err
		}
		close(ui.Error)
	}()

	ui.WaitInit()
}

func (ui *Control) AddTab(name string, action initAction) {
	lib.MainLogger.Println("add tab", name)
	ui.state.dirty()
	ui.tabs.add(name, action)
}

func (ui *Control) MarkTabError(name string, hasError bool) error {
	index, err := ui.tabs.findIndex(name)
	if err != nil {
		return err
	}
	ui.state.setErrorStatus(index, hasError)
	return nil
}

func (ui *Control) Update(view *gocui.View) {
	// lib.MainLogger.Println("try update view:", view.Name())
	i, _ := ui.tabs.findIndex(view.Name())
	if i == -1 || i == ui.tabs.current {
		ui.gui.Update(emptyRenderEvent)
		// lib.MainLogger.Println("update view:", view.Name())
	}
	return
}
