package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"yamato-sorariku/pubsub-emulator-ui/src/controller"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
		},
		// 許可したいHTTPリクエストヘッダの一覧
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"*",
		},
		MaxAge: 24 * time.Hour,
	}))

	// API namespace
	v1 := router.Group("/api/v1")
	{
		v1.POST("/pubsub/publish", controller.Publish)
	}

	router.GET("/ws", controller.HandleClients)

	go func() {
		for {
			_ = controller.PullPubSubMessage()
			time.Sleep(1 * time.Second)
		}
	}()

	controller.SetUpPubSub()

	router.Run(":8080")

}
