package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"test/pkg/db/models"
)

var (
	DBName         string = "ideanest"
	CollectionName string = "users"
	client         *mongo.Client
	collection     *mongo.Collection
)

func InitDB(DBName string, CollectionName string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	collection = client.Database(DBName).Collection(CollectionName)
	return nil
}
func CloseDB() error {
	return client.Disconnect(context.Background())
}

func GetAllUsers() ([]models.User, error) {

	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return nil, err
	}

	// make sure to close db after this function complete
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	// Define a filter to be used in the MongoDB query (empty to select all)
	filter := bson.M{}

	//execute mongodb operation and get a cursor to parse results
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	// make sure to close the cursor after this function complete
	defer cursor.Close(context.Background())

	// a slice to get all documents from mongodb user collection
	var users []models.User
	// cursor.All method, which is used to decode all documents from the MongoDB cursor into a slice
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}
func GetUserById(userID string) (*models.User, error) {
	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return nil, err
	}

	// make sure to close db after this function completes
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	// cast string to mongo obj id type

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	// Define a filter to find a user by ID
	filter := bson.M{"_id": objID}
	// Execute MongoDB operation to find a user by ID
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err
}
func GetUserByEmail(userEmail string) (*models.User, error) {
	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return nil, err
	}

	// make sure to close db after this function completes
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	// Define a filter to find a user by email
	filter := bson.M{"email": userEmail}

	// Execute MongoDB operation to find a user by email
	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, errors.New("Wrong Email or Password")
	}

	return &user, nil
}
func LoginUser(userEmail string, password string) (*models.User, error) {
	var user *models.User
	user, err := GetUserByEmail(userEmail)
	if err != nil {
		return nil, err
	}
	if user.Password != password {

		return nil, errors.New("Wrong Email or Password")
	}
	return user, nil
}
func UpdateUser(UserID string, updatedUser models.User) error {
	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return err
	}

	// Make sure to close the database after this function completes
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	objID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}

	// Execute MongoDB operation to update the user
	update := bson.M{"$set": bson.M{"name": updatedUser.Name, "email": updatedUser.Email, "password": updatedUser.Password}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	return nil
}
func InsertUser(newUser models.UserBind) error {
	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return err
	}

	// Make sure to close the database after this function completes
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	// Check if the email already exists
	existingUser := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		// Email already exists, return an error
		return fmt.Errorf("email %s already exists", newUser.Email)
	} else if err != mongo.ErrNoDocuments {
		// Other error occurred
		return err
	}

	// Execute MongoDB operation to insert a new user
	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return err
	}

	return nil
}
func DeleteUserById(userID string) error {
	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return err
	}

	// make sure to close db after this function completes
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	// Define a filter to find a user by ID
	filter := bson.M{"_id": objID}

	// Execute MongoDB operation to delete a user by ID
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
