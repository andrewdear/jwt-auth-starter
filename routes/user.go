package routes

import (
	"jwt-auth-starter/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AttachUserRoutes(router *gin.Engine) {
	router.POST("/signup", signUp)
	router.POST("/signin", signIn)
	router.GET("/profile", auth.RequiresAuth, getProfile)
}

func signUp(c *gin.Context) {
	//Get users credentials from the request
	user, suErr := auth.SignUp(c)

	if suErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": suErr.Error(),
		})
		return
	}
	// after signup create jwt token for user
	jwtToken, err := auth.GenerateJWT(user)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "error generating token",
		})
		return
	}
	// send back token to user
	c.JSON(200, gin.H{
		"token": jwtToken,
	})
}

func signIn(c *gin.Context) {

	//Get users credentials from the request
	user, siErr := auth.SignIn(c)

	if siErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect Credentials",
		})
		return
	}

	jwtToken, err := auth.GenerateJWT(user)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "error generating token",
		})
		return
	}

	c.JSON(200, gin.H{
		"token": jwtToken,
	})

}

func getProfile(c *gin.Context) {
	err, user := auth.GetUser(c)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, user)
}
