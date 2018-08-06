package util

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
