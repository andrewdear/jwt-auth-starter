package auth

import (
	"jwt-auth-starter/services"
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func createTestContext(body string, token string) *gin.Context {
	w := httptest.NewRecorder()

	request := httptest.NewRequest("POST","/test",
		bytes.NewBufferString(body))

	request.Header.Add("Token", token)
	request.Header.Add("Content-Type", "application/json")

	// create a testContext and set its request to one created above so i can put what i need to test with
	c, _ := gin.CreateTestContext(w)
	c.Request = request

	return c

}

var testHex string
var testToken string

//Test to make sure that we get back a user ID hex
func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := services.ConnectToMongo("jwt-auth-starter-test")

	if err != nil {
		t.Error("could not connect to database")
	}

	context := createTestContext(`{"email" : "testauth@test.com","password": "testPass"}`, "")

	hex, err := SignUp(context)

	if err != nil {
		t.Error(err)
	}

	if hex == "" {
		t.Error("Hex not returned")
	}

	testHex = hex

	// test it will not allow the same email to sign up two times

	hex, err = SignUp(context)

	if err == nil {
		t.Error("Sign up should have return error when using an email that exists")
	}

	if hex != "" {
		t.Error("Hex was returned when it should not have")
	}



}

func TestSignIn(t *testing.T) {

	context := createTestContext(`{"email" : "testauth@test.com","password": "testPass"}`, "")

	hex, err := SignIn(context)

	if err != nil {
		t.Error(err)
	}

	if hex == "" {
		t.Error("Hex not returned")
	}

}

// test that is generates a string and does not throw any errors
func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT(testHex)

	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token not returned")
	}

	testToken = token
}

func TestIsAuthorized(t *testing.T) {
	// Test if a correct token is passed in we get a hex

	ok, user := isAuthorized(testToken)

	if !ok {
		t.Error("user should be authorized but is not")
	}

	if user == "" {
		t.Error("no hex returned")
	}

	if user != testHex {
		t.Error("returned hex does not match the one passed in within generateJWT")
	}

}

// context to share between middleware tests and get user tests


func TestGinAuthMiddleWare(t *testing.T) {
	var context = createTestContext(``, testToken)

	GinAuthMiddleWare(context)
	//The middleware should take the token from the headers, resolve it and put it into the context
	// check that it has put it in the headers

	user, exists := context.Get("user")

	if !exists {
		t.Error("user was not put onto the context")
	}

	if user != testHex {
		t.Error("added the wrong hex to the context")
	}

}

func TestGetUserId(t *testing.T) {
	var context = createTestContext(``, testToken)
	// use the middleware to set the token
	GinAuthMiddleWare(context)
	// then make sure that the getUserId can get the user hex from the context
	user := GetUserId(context)

	if user != testHex {
		t.Error("unsuccessfully retrieved the user id from context")
	}

}

func TestGetUser(t *testing.T) {

	var context = createTestContext(``, testToken)
	// use the middleware to set the token
	GinAuthMiddleWare(context)

	// make sure we get a user back from getUser

	err, user := GetUser(context)

	if err != nil {
		t.Error(err)
	}

	if user.Email != "testauth@test.com" {
		t.Errorf("expected email to equal testauth@test.com but received %v", user.Email)
	}

	user.Delete()

}

// RequiredAuth is tested in routes