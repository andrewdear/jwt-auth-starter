package routes

import (
	"jwt-auth-starter/auth"
	"jwt-auth-starter/services"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter(database string) *gin.Engine {
	router := gin.Default()
	// connect to mongo database and reference to connection into services/index.go
	err := services.ConnectToMongo(database)

	if err != nil {
		log.Fatal(err)
	}

	// add in the jwt checker to each request
	router.Use(auth.GinAuthMiddleWare)

	AttachUserRoutes(router)

	return router
}