package serverbase

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type BaseServer struct {
	port   string
	router *gin.Engine
}

func NewBaseServer(port string) *BaseServer {
	base := new(BaseServer)
	base.port = port
	base.router = gin.Default()
	return base
}

func (b *BaseServer) AddRouter(path string, funcall gin.HandlerFunc) {
	b.router.GET(path, funcall)
}

func (b *BaseServer) StartServer() {
	if err := b.router.Run(":" + b.port); err != nil {
		fmt.Println("StartServer:", err)
	}
}
