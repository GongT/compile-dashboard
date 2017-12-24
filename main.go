package main

import (
	"github.com/gongt/compile-dashboard/lib/gui"
	"github.com/gongt/compile-dashboard/lib"
	"os"
	"os/signal"
	"syscall"
	"github.com/gongt/compile-dashboard/lib/process"
	"fmt"
)

func main() {
	lib.InitLogger()
	defer func() { lib.CloseLogger() }()

	ret := 0
	ch := createWatcher()

	screen, err := gui.NewControl()
	if err != nil {
		panic(err)
	}
	defer func() { screen.Close() }()

	go func() {
		<-ch
		lib.MainLogger.Println("receive quit signal")
		os.Exit(ret)
	}()

	// decodeConfigFile()
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

	lib.MainLogger.Println("main event loop started")
	if err := screen.EventLoop(); err != nil {
		panic(err)
	}
	//process1.Close()
	lib.MainLogger.Println("main app terminated")
}

/*
defer func () {
	if r := recover(); r != nil {
		screen.Close()
		os.Exit(1)
	}
}()
*/
func createWatcher() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGQUIT)
	return ch
}
