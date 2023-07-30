package main

import (
	"github.com/gin-gonic/gin"
	"go-basic/webbook/internal/web"
	"log"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	log.Fatal(server.Run(":8080"))
}
