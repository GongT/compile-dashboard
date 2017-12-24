package gui

import "github.com/jroimartin/gocui"

const viewNameAll, viewNameStatus, viewNameMessage, viewNameBottom = "", "status", "message", "bottom"
const maxTabCount = 20

type initAction func(event ViewInitEvent)

type ViewInitEvent struct {
	Gui  *gocui.Gui
	View *gocui.View
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

func emptyRenderEvent(gui *gocui.Gui) error {
	return nil
}
