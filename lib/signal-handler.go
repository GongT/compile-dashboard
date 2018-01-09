package lib

import (
	"os"
	"os/signal"
	"syscall"
	"runtime/debug"
	"fmt"
	"sync"
)

var MainDebugHang chan int

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
var mutex sync.Mutex

func init() {
	mutex = sync.Mutex{}
	runningThing = []*Stoppable{}
	MainDebugHang = make(chan int)
}

type UnregisterCallback = func()

func RegisterRunning(thing Stoppable) UnregisterCallback {
	runningThing = append(runningThing, &thing)
	return func() {
		UnRegisterRunning(thing)
		thing.Close()
	}
}
func UnRegisterRunning(thing Stoppable) {
	mutex.Lock()
	for index, item := range runningThing {
		if item == &thing {
			runningThing = append(runningThing[:index], runningThing[index+1:]...)
			break
		}
	}
	mutex.Unlock()
}

func StartSignalHandler() {
	ch := createWatcher()

	go func() {
		MainLogger.Println("receive quit signal:", <-ch)
		GracefulQuit(0)
	}()
}

func safeCleanup(obj *Stoppable) {
	defer func() {
		if r := recover(); r != nil {
			MainLogger.Printf("failed cleanup: %s\n%s", r, debug.Stack())
		}
	}()
	MainLogger.Println("stopping", (*obj).Inspect())
	fmt.Println("stopping", (*obj).Inspect())
	(*obj).Close()
}

func GracefulQuit(exit int) {
	println("\x1B[H\x1B[JNotice: will terminate all child process, please wait...")
	MainLogger.Println("GracefulQuit: ", exit)
	for index := len(runningThing) - 1; index >= 0; index-- {
		safeCleanup(runningThing[index])
	}
	MainLogger.Println("os.Exit() call")
	defer func() {
		println("bye~")
		os.Exit(exit)
	}()

	close(MainDebugHang)
}
