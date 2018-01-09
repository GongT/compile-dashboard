package gui

type renderEvent struct {
	ctl             *Control
	splitSideMiddle int
	splitCenter     int
	maxX            int
	maxY            int
}

type renderable interface {
	render(r *renderEvent)
}

func newRenderEvent(gui *Control) *renderEvent {
	maxX, maxY := gui.gui.Size()

	return &renderEvent{
		gui,
		int(float32(maxY-5) / 2),
		int(0.8 * float32(maxX)),
		maxX,
		maxY,
	}
}
