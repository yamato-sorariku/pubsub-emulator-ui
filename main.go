package main

import (
	"github.com/gin-gonic/gin"
	"yamato-sorariku/gcp-pubsub-ui/src/controller"
)

func main() {
	router := gin.Default()

	// API namespace
	v1 := router.Group("/api/v1")
	{
		v1.POST("/pubsub", controller.Pubsub)
		v1.GET("/pubsub/setEndPoint", controller.SetPushEndpoint)
		v1.GET("/pubsub/publish", controller.Publish)
	}

	router.Run(":8080")
}
