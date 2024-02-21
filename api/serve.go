package api

import (
	"palco-planner-api/controller"

	"github.com/gin-gonic/gin"
)

type Server struct {
	controller controller.Controller
}

func (s Server) Start() {
	router := gin.Default()

	router.POST("upload", s.controller.Upload)
	router.GET("files", s.controller.ListFilesInFolder)

	router.Run(":8080")
}

func NewServer(controller controller.Controller) Server {
	return Server{
		controller,
	}
}
