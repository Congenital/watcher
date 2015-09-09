package main

import (
	"github.com/Congenital/log/v0.2/log"
	"github.com/Congenital/watcher"
	"os"
	"runtime"
	"syscall"
	"time"
)

func main() {
	go func() {
		for {

			<-time.Tick(time.Second)
			log.Info(runtime.NumGoroutine())
		}
	}()

	watcher.SetBuffSize(10)
	watcher.Listen()

	watcher.SetHandle(syscall.SIGHUP, func() {
		log.Info("SIGHUP")
	})

	watcher.SetHandle(syscall.SIGINT, func() {
		log.Info("SIGINT")
		watcher.Exit(0)
	})

	watcher.SetDefaultHandle(func(_signal os.Signal) {
		log.Info("Default ---", _signal)
		watcher.ClearDefaultHandle()
	})

	watcher.GetExitCode()

}
