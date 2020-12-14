package handler

import (
	"log"

	"github.com/gin-gonic/gin"

	"restaurant-assistant/pkg/app/server"
)

func Init(server *server.RestaurantAssistantServer) {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/auth/login", server.Login)
		apiV1.POST("/auth/token/refresh", AuthMiddleware(), server.RefreshToken)
		apiV1.DELETE("/auth/logout", AuthMiddleware(), server.Logout)

		apiV1.POST("/user/create", server.CreateUser)
		apiV1.PUT("/user/create", AuthMiddleware(), server.UpdateUser)

		apiV1.POST("/product/create", AuthMiddleware(), server.CreateProduct)
		apiV1.POST("/order/create", AuthMiddleware(), server.CreateOrder)
	}

	log.Fatal(router.Run(":8080"))
}
