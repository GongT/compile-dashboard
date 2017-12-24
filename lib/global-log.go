package lib

import (
	"os"
	"log"
)

var MainLogger *log.Logger
var f *os.File

func InitLogger() {
	f, err := os.Create(os.TempDir() + "/compile-dashboard-temp.log")
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
