package server

import (
	"fmt"
	"gAD-System/services/gad-manager/domain"
	"gAD-System/services/gad-manager/models"

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
	router.POST("/calc", h.calculate)

	return router
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h Handlers) calculate(c *gin.Context) {
	var reqBody models.CalcRequest
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(reqBody)

	ans, err := h.Calculator.Calculate(reqBody.Exprs)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"ans": ans,
	})
}
