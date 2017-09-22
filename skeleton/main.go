package skeleton

import (
	"fmt"
	"store/core"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run(*core.Core)
	Shutdown()
	HttpServerHanlde(c *gin.Context)
}

type App struct {
	s    Server
	name string
	path string
	need *[]string
}

var appArray []*App

func AddServer(name string, s Server, path string, need *[]string) error {

	for _, app := range appArray {
		if app.name == name {
			return fmt.Errorf("Server %v already Exist")
		}
	}

	app := &App{
		name: name,
		s:    s,
		path: path,
		need: need,
	}

	appArray = append(appArray, app)

	return nil
}

func RunServers(coreData *core.Core) {

	var wg sync.WaitGroup

	runflags := make(map[string]bool, 0)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	for _, s := range appArray {

		if s.need != nil {
			for _, need := range *s.need {
				if runflags[need] == false {
					fmt.Println("Server ", s.name, " need server", need, " run first!")
					return
				}
			}

		}
		runflags[s.name] = true

		fmt.Println("Server Run :", s.name, s.path)

		wg.Add(1)
		go func(s *App) {
			router.GET(s.path, s.s.HttpServerHanlde)
			s.s.Run(coreData)
			s.s.Shutdown()
			wg.Done()
		}(s)

	}

	go router.Run(":" + coreData.Cont.Port)

	wg.Wait()
}

func StopServers() {

	for i := len(appArray) - 1; i >= 0; i-- {
		appArray[i].s.Shutdown()
	}

}
