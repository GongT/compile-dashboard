package gui

import (
	"github.com/jroimartin/gocui"
)

const viewNameAll, viewNameStatus, viewNameMessage, viewNameBottom = "", "status", "message", "bottom"
const maxTabCount = 20

type initAction func(event ViewInitEvent)

type ViewInitEvent struct {
	gui  *Control
	view *gocui.View
}

func newViewInitEvent(gui *Control, view *gocui.View) ViewInitEvent {
	return ViewInitEvent{
		gui,
		view,
	}
}
func (event ViewInitEvent) Write(bytes []byte) (n int, err error) {
	n, err = event.view.Write(bytes)
	event.gui.Update(event.view)
	return
}
func (event ViewInitEvent) Clear() {
	event.view.Clear()
	event.gui.Update(event.view)
}


func subLayout(g *gocui.Gui, name string, geo Geo, init func(*gocui.View)) error {
	view, err := g.SetView(name, geo.left, geo.top, geo.right, geo.bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		init(view)
	}

	return nil
}

func findView(g *gocui.Gui, name string) *gocui.View {
	item, err := g.View(name)
	if err != nil {
		return nil
	}
	return item
}

type Geo struct {
	left, top, right, bottom int
}

func emptyRenderEvent(_ *gocui.Gui) error {
	return nil
}
