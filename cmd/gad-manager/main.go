package main

import (
	"gAD-System/services/gad-manager/config"
	"gAD-System/services/gad-manager/server"
)

func main() {
	cfg := config.InitConfig()
	server.InitREST(cfg)

	// defer rpcCalcConn.
	// r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // :8080
}
