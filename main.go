package main

import (
	"os"
	"os/signal"
	"runtime"
	"store/cmd"
	"syscall"
)

type IServer interface {
	OnStart()
	OnShutdown()
}

func interceptSignal(s IServer) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.OnShutdown()
		os.Exit(0)
	}()
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	server := new(cmd.StoreServer)

	interceptSignal(server)

	server.OnStart()

	return
}
