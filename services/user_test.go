package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

var savedUserObjectId primitive.ObjectID
var email = "goTest@goTest.com"

//TESTS NEED TO BE RUN TOGETHER AS THEY DEPENDS ON EACH OTHER AND THE CONNECTION SETUP IN SAVE

// Test to makes ure that User.Save, saves the user if it has a email and password attached
func TestUser_Save(t *testing.T) {
	err := ConnectToMongo("jwt-auth-starter-test")

	if err != nil {
		t.Error("could not connect to database")
	}

	user := User{Email: email, Password: "password"}

	err = user.Save()


	if err != nil {
		t.Error("could not save user", err)
	}

	if user.ID.IsZero() {
		t.Error("could not save user", err)
	}

	// save object id to global var so can be used in other tests.
	savedUserObjectId = user.ID


}

func TestUser_Get(t *testing.T) {
	// test get user by Id
	user := User{ID: savedUserObjectId}

	err := user.Get()

	if err != nil {
		t.Error(err)
	}

	if user.Email != email {
		t.Errorf("expected email to equal %v but received %v", email, user.Email)
	}

	// Test get user by Email

	user2 := User{Email: email}

	err = user2.Get()

	if err != nil {
		t.Error(err)
	}

	if user.ID != savedUserObjectId {
		t.Errorf("expected ID to equal %v but received %v when", savedUserObjectId, user.ID)
	}

}

func TestUser_Save_updateMode(t *testing.T) {
	// test that a user can be update with save,
	user := User{ID: savedUserObjectId}

	user.Get()

	testDescription := "Test description"

	user.Description = testDescription

	err := user.Save()

	if err != nil {
		t.Error(err)
	}

	// re retreive the user from a fresh instance and make sure that the updated description is there
	user = User{ID: savedUserObjectId}

	user.Get()

	if user.Description != testDescription {
		t.Error("description has not been updated")
	}

}

func TestUser_CheckIfExists(t *testing.T) {
	// check that a user that is empty returns false
	user := User{}

	err, exists := user.CheckIfExists()

	if err == nil {
		t.Error("should throw an error if no information within user")
	}

	if exists {
		t.Error("empty user struct should not exist")
	}

	// check that a user with an ID field returns that it exists
	user = User{ID: savedUserObjectId}

	err, exists = user.CheckIfExists()

	if !exists {
		t.Error("Expected to find the user with id")
	}

	// check that a user with an email field that exists returns as true
	user = User{Email: email}

	err, exists = user.CheckIfExists()

	if !exists {
		t.Error("Expected to find the user with email")
	}

	// check that a user with an bad email field returns as false
	user = User{Email: "doesnotexist@bad.com"}

	err, exists = user.CheckIfExists()

	if exists {
		t.Error("Expected user with bad email to not exist")
	}

}

func TestUser_Delete(t *testing.T) {
	user := User{ID: savedUserObjectId}

	err := user.Delete()

	if err != nil {
		t.Error(err)
	}

	if !user.ID.IsZero() {
		t.Error("user ID should be zero after the user has been deleted")
	}

}