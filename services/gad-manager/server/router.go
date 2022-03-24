package server

import (
	"gAD-System/services/gad-manager/domain"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	calculator domain.Calculator
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// router.GET("/health", health.Status)
	// router.Use(middlewares.AuthMiddleware())

	router.GET("/ping", Pong)

	return router
}

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h Handlers) Calculate(c *gin.Context) {
	ans, err := h.calculator.Calculate("100+100")
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"ans": ans,
	})
}
