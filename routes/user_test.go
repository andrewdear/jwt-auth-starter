package routes

import (
	"jwt-auth-starter/services"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method string, target string, body string, token string) *httptest.ResponseRecorder {
	// we create a test request
	request := httptest.NewRequest(method,target, bytes.NewBufferString(body))
	request.Header.Add("Content-Type", "application/json")
	if token != "" {
		request.Header.Add("Token", token)
	}

	// create a recorder to save the request data
	w := httptest.NewRecorder()
	// server up a gin http with the writer and request
	r.ServeHTTP(w, request)
	return w
}

var testToken string

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter("jwt-auth-starter-test")
	// test request to Signup
	w := performRequest(router,"POST", "/signup", `{"email" : "testroutes@test.com","password": "testPass"}`, "")

	if w.Code != 200 {
		t.Errorf("expected a status of 200 but received %v", w.Code)
	}

	var response map[string]string

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	if err != nil {
		t.Error("error decoding the json request")
	}

	// check that a token has been returned in the response
	if val, ok := response["token"]; ok {
		if !ok {
			t.Error("no token in response")
		}

		if val == "" {
			t.Error("token empty in response")
		}

	}

	// Check that it will not allow the same email to signup two times

	w = performRequest(router,"POST", "/signup", `{"email" : "testroutes@test.com","password": "testPass"}`, "")

	if w.Code != 401 {
		t.Errorf("expected a status of 401 but received %v", w.Code)
	}

	var badResponse map[string]string

	err = json.Unmarshal([]byte(w.Body.String()), &badResponse)

	if err != nil {
		t.Error("error decoding the json request")
	}

	// check that a token has been returned in the badResponse
	if val, ok := badResponse["error"]; ok {
		if !ok {
			t.Error("no error in badResponse")
		}

		if val != "email already in use" {
			t.Errorf("incorrect error expected 'email already in use' received %v", val)
		}

	}

}

func TestLogIn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter("jwt-auth-starter-test")
	// test request to Signin with correct credentials
	w := performRequest(router,"POST", "/signin", `{"email" : "testroutes@test.com","password": "testPass"}`, "")

	if w.Code != 200 {
		t.Errorf("expected a status of 200 but received %v", w.Code)
	}

	var response map[string]string

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	if err != nil {
		t.Error("error decoding the json request")
	}

	// check that a token has been returned in the response
	if val, ok := response["token"]; ok {
		if !ok {
			t.Error("no token in response")
		}

		if val == "" {
			t.Error("token empty in response")
		}
		// save the token for use the profile test
		testToken = val

	}

	// ensure it will not allow you to login with incorrect credentials

	w = performRequest(router,"POST", "/signin", `{"email" : "testroutes@test.com","password": "testPa"}`, "")

	if w.Code != 401 {
		t.Errorf("expected a status of 401 but received %v", w.Code)
	}

	var bedResponse map[string]string

	err = json.Unmarshal([]byte(w.Body.String()), &bedResponse)

	if err != nil {
		t.Error("error decoding the json request")
	}

	// check that a token has been returned in the badResponse
	if val, ok := bedResponse["error"]; ok {
		if !ok {
			t.Error("no error in badResponse")
		}

		if val != "Incorrect Credentials" {
			t.Errorf("incorrect error expected 'Incorrect Credentials' received %v", val)
		}

	}


}

func TestGetProfile(t *testing.T) {
	//TEST THAT A CORRECT TOKEN WILL RETURN USER DETAILS
	gin.SetMode(gin.TestMode)
	router := SetupRouter("jwt-auth-starter-test")
	// test request to Signin with correct credentials
	w := performRequest(router,"GET", "/profile", ``, testToken)

	if w.Code != 200 {
		t.Errorf("expected a status of 200 but received %v", w.Code)
	}

	var response map[string]string

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	if err != nil {
		t.Error("error decoding the json request")
	}

	// check that a token has been returned in the response
	if val, ok := response["email"]; ok {
		if !ok {
			t.Error("no email in response")
		}

		if val != "testroutes@test.com" {
			t.Errorf("expected email to be testroutes@test.com but received %v", val)
		}

	}

	//TEST THAT A INCORRECT TOKEN WILL RETURN AN ERROR

	w = performRequest(router,"GET", "/profile", ``, "testToken")

	if w.Code != 401 {
		t.Errorf("expected a status of 401 but received %v", w.Code)
	}

	var bedResponse map[string]string

	err = json.Unmarshal([]byte(w.Body.String()), &bedResponse)

	if err != nil {
		t.Error("error decoding the json request")
	}

	// check that a token has been returned in the badResponse
	if val, ok := bedResponse["error"]; ok {
		if !ok {
			t.Error("no error in badResponse")
		}

		if val != "Unauthorized" {
			t.Errorf("incorrect error expected 'Unauthorized' received %v", val)
		}

	}

	// use the results from the good request to delete the test user
	userObjectId, _ := primitive.ObjectIDFromHex(response["id"])

	user := services.User{ID: userObjectId}

	user.Delete()
}