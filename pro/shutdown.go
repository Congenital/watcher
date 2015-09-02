package main

import (
	"github.com/Congenital/log/v0.2/log"
	"github.com/Congenital/watcher"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Info("Please input process id :")
		return
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Error(err)
		return
	}

	watcher.ShutDown(pid)
}
