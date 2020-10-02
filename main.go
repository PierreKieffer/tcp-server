package main

import (
	"github.com/PierreKieffer/tcp-server/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	srv := server.InitServer()

	select {
	case <-exit:
		srv.Stop()
	}
}
