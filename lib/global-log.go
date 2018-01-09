package lib

import (
	"os"
	"log"
	"fmt"
)

var MainLogger *log.Logger
var f *os.File

func println(args ... interface{}) {
	fmt.Fprintln(os.Stdout, args...)
}

func InitLogger() {
	f, err := os.OpenFile(os.TempDir()+"/compile-dashboard-temp.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	MainLogger = log.New(f, "", log.Lshortfile)
	MainLogger.Println("Hello!")
	f.Sync()
}

func CloseLogger() {
	f.Close()
}
