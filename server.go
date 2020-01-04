package main

import (
	"github.com/gin-gonic/gin"
	"gp-websoket/api"
	"gp-websoket/database"
	"gp-websoket/middleware"
)

const HOST = "127.0.0.1:8888"

func main() {
	database.Migrate()

	r := gin.Default()

	r.POST("/login", api.Login)
	authorized := r.Group("/").Use(middleware.JWTAuth)
	authorized.GET("/rooms", api.Rooms)
	authorized.POST("/rooms", api.CreateRoom)
	authorized.POST("/rooms/:name/join", api.JoinRoom)
	authorized.GET("/ws", api.WsHandler)

	r.Run(HOST) // 启动服务并监听
}
