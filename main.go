package main

import (
	"github.com/gongt/compile-dashboard/lib/gui"
	"github.com/gongt/compile-dashboard/lib"
	"github.com/gongt/compile-dashboard/lib/process"
	"fmt"
	"github.com/jroimartin/gocui"
)

func main() {
	lib.InitLogger()
	defer func() { lib.CloseLogger() }()

	startSignalHandler()

	screen, err := gui.NewControl()
	if err != nil {
		panic(err)
	}
	registerRunning(screen)

	go func() {
		err := <-screen.Error
		lib.MainLogger.Println("screen control result:", err)
		if err == gocui.ErrQuit {
			GracefulQuit(0)
		} else {
			GracefulQuit(100)
		}
		lib.MainLogger.Println("program must quit now!")
	}()

	screen.StartEventLoop()

	screen.Message(fmt.Sprint("configFile:", lib.ConfigFile))

	// output := lib.NewPseudoOutput()
	// output.Start()

	/*screen.AddTab("Tab A", func(event gui.ViewInitEvent) {
		go func() {
			for {
				data := <-output.Tunnel
				fmt.Println(data)
				screen.OutputTab("Tab A", []byte(data))
			}
		}()
	})*/

	screen.AddTab("Tab A", func(event gui.ViewInitEvent) {
		lib.MainLogger.Println("Tab A has inited")
		view := event.View
		go func() {
			proc := process.NewChildProcess("while true; do echo -n 'process1:'; date; sleep 1; done")
			registerRunning(proc)
			defer func() {
				screen.Update(view)
				proc.Close()
			}()
			for {
				select {
				case data := <-proc.OutputPipe.Output:
					view.Write(data)
				case <-proc.OutputPipe.Clear:
					view.Clear()
				case err := <-proc.Stop:
					screen.MarkTabError("Tab A", true)
					fmt.Fprintln(view, "\n\nProcess Stoped: ", err)
					return
				}
				screen.Update(view)
			}
		}()
	})

	screen.AddTab("Tab B", func(event gui.ViewInitEvent) {
		view := event.View
		go func() {
			proc := process.NewChildProcess("truncate --size 0 test1.txt; tail -f test1.txt")
			registerRunning(proc)
			defer func() {
				screen.Update(view)
				proc.Close()
			}()
			for {
				select {
				case data := <-proc.OutputPipe.Output:
					view.Write(data)
					lib.MainLogger.Println("Tab B output!", string(data))
				case <-proc.OutputPipe.Clear:
					view.Clear()
					lib.MainLogger.Println("Tab B Clear!")
				case err := <-proc.Stop:
					screen.MarkTabError("Tab B", true)
					fmt.Fprintln(view, "\n\nProcess Stoped: ", err)
					lib.MainLogger.Println("Tab B Stoped!")
					return
				}
				screen.Update(view)
			}
		}()
	})

	<-mainDebugHang
}
