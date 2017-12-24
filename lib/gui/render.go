package gui

import "github.com/jroimartin/gocui"

type renderEvent struct {
	g               *gocui.Gui
	splitSideMiddle int
	splitCenter     int
	maxX            int
	maxY            int
}

type renderable interface {
	render(r *renderEvent)
}

func newRenderEvent(gui *gocui.Gui) *renderEvent {
	maxX, maxY := gui.Size()

	return &renderEvent{
		gui,
		int(float32(maxY-5) / 2),
		int(0.8 * float32(maxX)),
		maxX,
		maxY,
	}
}
