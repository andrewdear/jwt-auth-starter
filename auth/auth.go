package auth

import (
	"jwt-auth-starter/services"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

//TODO: https://godoc.org/github.com/gin-gonic/gin#CreateTestContext
// create test contet to pass to function, will need to figure out how to add params to it.
//https://github.com/gin-gonic/gin/blob/master/context_test.go  -- line 657to set context json body maybe?

type loginRequest struct {
	Email string `json:"email"`
	Password  string `json:"password"`
}

// you can set a token in the terminal by using set MY_JWT_SECRET = mySuperSecretPhrase
// then get it in go using
// var mySigningKey = []byte(os.Get("MY_JWT_SECRET"))

// the key should be retreived from an enviromental variable
var mySigningKey = []byte("mySuperSecretPhrase")

//SignUp gets the users crentials and then generates a hashed password then saves the user to the database
func SignUp(c *gin.Context) (string, error) {
	credentials := getCredentials(c)

	passwordByteSlice := []byte(credentials.Password)

	hashedPassword, errGenerate := bcrypt.GenerateFromPassword(passwordByteSlice, bcrypt.DefaultCost)
	if errGenerate != nil {
		return "", errGenerate
	}

	// create a new user
	user := services.User{
		Email: credentials.Email,
		Password: string(hashedPassword),
	}

	err, exists := user.CheckIfExists()

	if err != nil {
		return "", err
	}

	// if the exist check comes back as false then this is a new email so it is ok to create a new account
	if !exists {
		// save the user then get the id from it and return it
		err := user.Save()

		if err != nil {
			return "", nil
		}

		//Otherwise return the users id hex to put into a jwt cookie
		return user.ID.Hex(), nil
	}

	//otherwise return an error to send in the json
	return "" , errors.New("email already in use")

}

// SignIn gets the credentials from the request and compares that password to any saved to users credentials
func SignIn(c *gin.Context) (string, error) {
	credentials := getCredentials(c)

	// create a user struct so you can get the data from the database
	user := services.User{
		Email: credentials.Email,
	}

	err := user.Get()

	if err != nil {
		return "", err
	}

	// put the credentials password into the format required for compare
	passwordByteSlice := []byte(credentials.Password)

	// put the users hashed passworf from the database into the format required for compare
	hahsedPasswordByteSlice := []byte(user.Password)

	// passwordsMatchErr should be nil if passwords match
	passwordsMatchErr := bcrypt.CompareHashAndPassword(hahsedPasswordByteSlice, passwordByteSlice)

	if passwordsMatchErr != nil {
		return "" , errors.New("incorrect credentials")
	}


	// if password matches one in db for email then return the email to put into the jwt
	return user.ID.Hex() , nil
}

// GetCredentials takes credentials from request and puts them into a loginRequest
func getCredentials(c *gin.Context) *loginRequest {
	requestCredentials := loginRequest{}

	// this is saying give me a struct of what the request that will be sent in will look like and i will assign those values to that struct for you.
	c.Bind(&requestCredentials)

	return &requestCredentials

}

// GenerateFWT, this generates a JWT string to be passed to a client
func GenerateJWT(user string) (string, error) {
	// this defines the type of encoding we want the string to be in
	token := jwt.New(jwt.SigningMethodHS256)

	// claims are just peices of information that are in the token that can be used or retreived later.
	claims := token.Claims.(jwt.MapClaims)

	// need to pass these into function on create so can be taken from db
	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// we then get the token string and salt it with a secret key that needs to be a string byte slice
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func isAuthorized(token string) (bool, string) {
		//We get a token passed in to check

		if token != "" {
			// we then put a funcion into the parse to tell it how to parse it
			// the function get the token passed to it by the parse function.
			token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				// so this checks that the token can be parsed using the same parsing method from the client and if it is not ok we returns nothing
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return false, fmt.Errorf("there was an error")
				}

				// otherwise we return the signingKey from the parse function which should allow the parse function to return the token data
				return mySigningKey, nil
			})

			if err != nil {
				return false, ""
			}

			// if we have received the signing key from the token then we can send back the positive response
			if(token.Valid) {
				// this converts the token.Claims to a map that we can get the data we added into it from.
				return true, token.Claims.(jwt.MapClaims)["user"].(string)
			}

		}

		return false, ""


}


//GinAuthMiddleWare adds a middleware to Gin that looks at the headers for a Token, decodes it then sets the users id hex into the gin context
func GinAuthMiddleWare(c *gin.Context) {
	//Looks at the haeders and retreive the token if it is there
	token := c.Request.Header.Get("Token")

	ok, user := isAuthorized(token)

	// After checking to see if the user is authorized we put either the userId from the token or an empty string into the context for the routes to receive
	if !ok {
		c.Set("user", "")
	} else {
		c.Set("user", user)
	}

	c.Next()

}

//RequiresAuth: checks that a user if has been added to context from the middleware
func RequiresAuth(c *gin.Context) {

	user := GetUserId(c)

	if user == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()


}

//GetUserId return userId from the Gin context
func GetUserId(c *gin.Context) string {
	user, _ := c.Get("user")

	return user.(string)
}

//Get user uses the id to look up the user info and return it
func GetUser(c *gin.Context) (error, *services.User) {

		userHexId := GetUserId(c)

		userObjectId, _ := primitive.ObjectIDFromHex(userHexId)

		user := services.User{ID: userObjectId}

		err := user.Get()

		if err != nil {
			return err, nil
		}

		return nil, &user

}