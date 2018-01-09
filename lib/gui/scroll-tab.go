package gui

import (
	"github.com/jroimartin/gocui"
)

type scrollTab struct {
	name   string
	view   *gocui.View
	gui    *gocui.Gui
	action initAction
	inited bool
}

func newScrollTab(name string, gui *gocui.Gui, action initAction) *scrollTab {
	st := scrollTab{
		name,
		nil,
		gui,
		action,
		false,
	}

	err := st.initKeys()
	if err != nil {
		panic(err)
	}

	return &st
}

func (st *scrollTab) initKeys() error {
	if err := st.gui.SetKeybinding(st.name, gocui.MouseWheelUp, gocui.ModNone, func(_ *gocui.Gui, view *gocui.View) error {
		return st.mouseScrollUp(view)
	}); err != nil {
		return err
	}
	if err := st.gui.SetKeybinding(st.name, gocui.MouseWheelDown, gocui.ModNone, func(_ *gocui.Gui, view *gocui.View) error {
		return st.mouseScrollDown(view)
	}); err != nil {
		return err
	}
	return nil
}

func (st *scrollTab) render(e *renderEvent) error {
	err := subLayout(e.ctl.gui, st.name, Geo{0, 0, e.splitCenter - 1, e.maxY - 2}, func(view *gocui.View) {
		if st.inited {
			return
		}
		st.inited = true
		view.Title = "Application Output: " + st.name
		view.Wrap = true
		view.Autoscroll = true
		st.view = view
		go func() {
			st.action(newViewInitEvent(e.ctl, view))
			st.gui.Update(emptyRenderEvent)
		}()
	})
	if err != nil {
		return err
	}
	e.ctl.gui.SetViewOnTop(st.name)
	return nil
}

func (st *scrollTab) mouseScrollDown(view *gocui.View) error {
	x, y := view.Origin()

	_, height := view.Size()
	_, err := view.Line(y + 1 + height)
	if err != nil {
		if !view.Autoscroll {
			view.Autoscroll = true
		}
		return nil
	}

	view.SetOrigin(x, y+1)
	view.Autoscroll = false
	return nil
}
func (st *scrollTab) mouseScrollUp(view *gocui.View) error {
	x, y := view.Origin()
	view.SetOrigin(x, y-1)

	if view.Autoscroll {
		view.Autoscroll = false
	}

	return nil
}
