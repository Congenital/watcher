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
	w := watcher.NewWatcher()

	w.SetBuffSize(10)
	w.Listen()

	w.SetHandle(syscall.SIGHUP, func() {
		log.Info("SIGHUP")
	})

	w.SetHandle(syscall.SIGINT, func() {
		log.Info("SIGINT")
		w.Exit(0)
	})

	w.SetHandle(syscall.SIGUSR1, func() {
		log.Info("RESTART")
		w.Exit(1)
	})

	w.SetDefaultHandle(func(_signal os.Signal) {
		log.Info("Default ---", _signal)
		w.ClearDefaultHandle()
	})

	w.GetExitCode()

}
