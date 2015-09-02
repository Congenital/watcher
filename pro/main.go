package main

import (
	"github.com/Congenital/log/v0.2/log"
	"github.com/Congenital/watcher"
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
		watcher.Exit()
	})

	watcher.GetExitCode()

}
