package server

import (
	"gAD-System/services/gad-manager/domain"

	"github.com/gin-gonic/gin"
)

// Handlers stores dependecies for handling requests
type Handlers struct {
	Calculator *domain.Calculator
}

// InitServerHandlers inits handlers dependecies with upper-level entites
func InitServerHandlers(calculator *domain.Calculator) *Handlers {
	return &Handlers{Calculator: calculator}
}

func newRouter(h *Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", pong)
	router.GET("/calc", h.calculate)

	return router
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h Handlers) calculate(c *gin.Context) {
	ans, err := h.Calculator.Calculate([]string{"100+100", "200-20"})
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"ans": ans,
	})
}
