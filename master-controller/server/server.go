package server

import (
	"gAD-System/master-controller/config"

	"github.com/gin-gonic/gin"
)

func Init() {
	cfg := config.InitConfig()
	r := NewRouter()
	r.Run(cfg.REST.PortREST) // TODO error handle
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// router.GET("/health", health.Status)
	// router.Use(middlewares.AuthMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
