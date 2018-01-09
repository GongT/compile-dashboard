package main

import (
	"github.com/gongt/compile-dashboard/lib/gui"
	"github.com/gongt/compile-dashboard/lib"
	"github.com/gongt/compile-dashboard/lib/process"
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

func main() {
	lib.InitLogger()
	defer func() { lib.CloseLogger() }()

	lib.StartSignalHandler()

	lib.LoadConfigFile()

	screen, err := gui.NewControl()
	if err != nil {
		panic(err)
	}
	defer lib.RegisterRunning(screen)()

	go func() {
		err := <-screen.Error
		lib.MainLogger.Println("screen control result:", err)
		if err == gocui.ErrQuit {
			lib.GracefulQuit(0)
		} else {
			lib.GracefulQuit(100)
		}
		lib.MainLogger.Println("program must quit now!")
	}()

	screen.StartEventLoop()

	screen.Message(fmt.Sprint("configFile:", lib.ConfigFile))

	startMainTab(screen, lib.ConfigFile.Scripts.Start)

	<-lib.MainDebugHang
}

func startProcessorTab(screen *gui.Control, config lib.ConfigFileTranspilerDefine) {
	screen.AddTab(config.Title, func(view gui.ViewInitEvent) {
		proc := process.NewChildProcess("truncate --size 0 test1.txt; tail -f test1.txt")
		defer lib.RegisterRunning(proc)()
		for {
			select {
			case data := <-proc.OutputPipe.Output:
				view.Write(data)
			case <-proc.OutputPipe.Clear:
				view.Clear()
			case err := <-proc.Stop:
				if err != nil {
					screen.MarkTabError(config.Title, true)
				}
				fmt.Fprintln(view, "\n\nProcess Stoped: ", err)
				lib.MainLogger.Println(config.Title + " Stoped!")
				return
			}
		}
	})
}

func startMainTab(screen *gui.Control, mainScript string) {
	const name = "SERVER"
	screen.AddTab(name, func(view gui.ViewInitEvent) {
		proc := process.NewChildProcess(mainScript)
		lib.RegisterRunning(proc)

		defer func() { lib.UnRegisterRunning(proc) }()

		for {
			select {
			case data := <-proc.OutputPipe.Output:
				view.Write(data)
			case <-proc.OutputPipe.Clear:
				view.Clear()
			case err := <-proc.Stop:
				if err != nil {
					screen.MarkTabError(name, true)
				}
				fmt.Fprintln(view, "\n\nProcess Stoped: ", err)

				lib.MainLogger.Println(name + " Stoped!")

				// auto restart after 5s
				<-time.After(5 * time.Second)
				proc = process.NewChildProcess(mainScript)
			}
		}
	})
}
