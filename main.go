package main

import (
	"jwt-auth-starter/routes"
)

func main() {

	router := routes.SetupRouter("jwt-auth-starter")

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
