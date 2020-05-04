package services

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email string `bson:"email,omitempty" json:"email"`
	Password string `bson:"password,omitempty" json:"-"`
	Description string `bson:"description,omitempty" json:"description"`
}

// Save takes the user instance and creates or updates it in the database
func (user *User) Save() error {
	// because the get collection in in the index Service i can use it without import
	userCollection := GetCollection("user")

	err, exists := user.CheckIfExists()

	if err != nil {
		return nil
	}

	ctx, cancel := CreateContext()
	defer cancel()

	// if the user already exists and has and id so has been saved before, we need to update it
	if exists && !user.ID.IsZero() {

		// we dont care about the result we only care if it errored as it could not update if we try and save the same data
		_ , err := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.D{
				{"$set", user},
			},

		)

		if err != nil {
			return err
		}

		return nil

	}

	userResult, err := userCollection.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	user.ID = userResult.InsertedID.(primitive.ObjectID)

	return nil
}

// Get uses relevant fields of the user struct to find the user from the database and assign it back to the user.
func (user *User) Get() error {
	userCollection := GetCollection("user")

	// create user instance to save any data from database
	var foundUser User

	var err error

	ctx, cancel := CreateContext()
	defer cancel()

	//If we have a user id use that to find the user
	if !user.ID.IsZero()  {
		if err = userCollection.FindOne(ctx, bson.M{"_id":user.ID}).Decode(&foundUser); err != nil {
			return errors.New("user not found")
		}

		*user = foundUser
		return nil
	}

	// otherwise use the email to find the user
	if user.Email != "" {
		ctx, cancel := CreateContext()
		defer cancel()

		if err = userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser); err != nil {
			return errors.New("user not found")
		}

		*user = foundUser
		return nil
	}

	// if we dont have an id or email return error
	return errors.New("user not found")

}

func (user *User) Delete() error {
	userCollection := GetCollection("user")

	ctx, cancel := CreateContext()
	defer cancel()

	deleteResult, err := userCollection.DeleteOne(ctx, user)

	if err != nil {
		return err
	}
	if deleteResult.DeletedCount != 1 {
		return errors.New("delete count is not 1")
	}

	// after delete set the user to an empty user object
	*user = User{}

	return nil

}


//CheckIfExists looks to see if there is an ID or if not looks into the database using email to see if user exists
func (user *User) CheckIfExists() (error, bool) {

	// IF no information to do a query exits then return err and false
	if user.Email == "" && user.ID.IsZero() {
		return errors.New("no defining infomration in credentials"), false
	}

	userCollection := GetCollection("user")

	// create user instance to save any data from database
	var foundUser User

	// we only want to query on either id or email so make bson map with only those fields
	// we only want to query on either id or email so make bson map with only those fields
	searchQuery := bson.M{}

	if !user.ID.IsZero() {
		searchQuery["_id"] = user.ID
	}

	if user.Email != "" {
		searchQuery["email"] = user.Email
	}

	ctx, cancel := CreateContext()
	defer cancel()

	err := userCollection.FindOne(ctx, searchQuery).Decode(&foundUser)

	if err != nil {
		return nil, false
	}

	return nil, true

}