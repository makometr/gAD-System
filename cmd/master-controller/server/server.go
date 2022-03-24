package server

import (
	"gAD-System/master-controller/config"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init() {
	cfg := config.InitConfig()
	r := NewRouter()
	r.Run(cfg.REST.Port) // TODO error handle
}

func InitCalculateRPC(cfg *config.Config) *grpc.ClientConn {
	conn, err := grpc.Dial(cfg.RPCCalc.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
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
