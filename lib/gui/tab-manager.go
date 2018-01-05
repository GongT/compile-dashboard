package gui

import (
	"github.com/jroimartin/gocui"
	"errors"
	"github.com/gongt/compile-dashboard/lib"
)

type tabManager struct {
	current  int
	length   int
	tabs     [maxTabCount]*scrollTab
	tabNames [maxTabCount]*string
	gui      *gocui.Gui
	notInit  []*scrollTab
}

func newTabManager(gui *gocui.Gui) *tabManager {
	return &tabManager{
		0,
		0,
		[maxTabCount]*scrollTab{},
		[maxTabCount]*string{},
		gui,
		[]*scrollTab{},
	}
}

func (mm *tabManager) switchTo(index int) error {
	if index < 0 || index >= maxTabCount {
		return errors.New("tab index out of range")
	}
	mm.current = index
	mm.gui.Update(emptyRenderEvent)
	return nil
}

func (mm *tabManager) add(name string, action initAction) {
	if mm.length >= maxTabCount {
		panic(errors.New("tab index out of range"))
	}

	tab := newScrollTab(name, mm.gui, action)

	mm.tabs[mm.length] = tab
	mm.tabNames[mm.length] = &name

	mm.notInit = append(mm.notInit, tab)

	mm.length++
}

func (mm *tabManager) render(e *renderEvent) error {
	if mm.length == 0 {
		return nil
	}
	for index, needInit := range mm.notInit {
		lib.MainLogger.Println("first time init view: ", needInit.name)
		if index != mm.current {
			needInit.render(e)
		}
	}
	mm.notInit = []*scrollTab{}

	return mm.tabs[mm.current].render(e)
}
func (mm *tabManager) findIndex(name string) (int, error) {
	// lib.MainLogger.Println("try find index: ", mm.current)
	for index := 0; index < mm.length; index++ {
		// lib.MainLogger.Println("try find index: ", mm.tabNames[index], *mm.tabNames[index])
		if name == *mm.tabNames[index] {
			return index, nil
		}
	}
	return -1, errors.New("can not find tab: " + name)
}
