package cmd

import (
	"fmt"
	"store/core"
	"store/initstore"
	"store/servercalculat"
	"store/serverweb"
	"store/skeleton"
	"store/util/log"
)

type StoreServer struct {
}

func (s *StoreServer) OnStart() {

	if _, err := log.NewLog("/Users/fanhongyue/Desktop/golang/src/log", "store", 0); err != nil {
		fmt.Println("failed to initLog, err %s:", err.Error())
		return
	}

	s.initServer()

	core := core.NewCore()

	if err := initstore.OnInitCore(core); err != nil {
		panic(err)
	}

	skeleton.RunServers(core)

	return
}

func (s *StoreServer) OnShutdown() {
	skeleton.StopServers()
}

func (s *StoreServer) initServer() {

	skeleton.AddServer("webserver", serverweb.NewWebServer(), "/webservers", nil)
	skeleton.AddServer("calculatserver", servercalculat.NewCalculatServer(), "/calculatservers", &[]string{"webserver"})
}
