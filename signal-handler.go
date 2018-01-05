package main

import (
	"github.com/gongt/compile-dashboard/lib"
	"os"
	"os/signal"
	"syscall"
	"runtime/debug"
	"fmt"
)

var mainDebugHang chan int

func createWatcher() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGQUIT)
	return ch
}

type Stoppable interface {
	Inspect() string
	Close() error
}

var runningThing []*Stoppable

func init() {
	runningThing = []*Stoppable{}
	mainDebugHang = make(chan int)
}

func registerRunning(thing Stoppable) {
	runningThing = append(runningThing, &thing)
}

func startSignalHandler() {
	ch := createWatcher()

	go func() {
		lib.MainLogger.Println("receive quit signal:", <-ch)
		GracefulQuit(0)
	}()
}

func safeCleanup(obj *Stoppable) {
	defer func() {
		if r := recover(); r != nil {
			lib.MainLogger.Printf("failed cleanup: %s\n%s", r, debug.Stack())
		}
	}()
	lib.MainLogger.Println("stopping", (*obj).Inspect())
	fmt.Println("stopping", (*obj).Inspect())
	(*obj).Close()
}

func GracefulQuit(exit int) {
	println("\x1B[H\x1B[JNotice: will terminate all child process, please wait...")
	lib.MainLogger.Println("GracefulQuit: ", exit)
	for index := len(runningThing) - 1; index >= 0; index-- {
		safeCleanup(runningThing[index])
	}
	lib.MainLogger.Println("os.Exit() call")
	defer func() {
		println("bye~")
		os.Exit(exit)
	}()

	close(mainDebugHang)
}
