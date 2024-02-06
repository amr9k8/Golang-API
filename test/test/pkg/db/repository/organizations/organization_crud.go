package organizations

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"test/pkg/db/models"
	"test/pkg/db/repository/users"
)

var (
	DBName         string = "ideanest"
	CollectionName string = "organizations"
	client         *mongo.Client
	collection     *mongo.Collection
)

func InitDB(DBName string, CollectionName string) error {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
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

func GetAllOrganizations() ([]models.Organization, error) {
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

	//execute MongoDB operation and get a cursor to parse results
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	// make sure to close the cursor after this function complete
	defer cursor.Close(context.Background())

	// a slice to get all documents from MongoDB organization collection
	var organizations []models.Organization
	// cursor.All(), a method, which is used to decode all documents from the MongoDB cursor into a slice
	if err := cursor.All(context.Background(), &organizations); err != nil {
		return nil, err
	}

	// Loop through each organization
	for _, org := range organizations {
		// Loop through each member_id in the organization
		for i := range org.Members {
			user, _ := users.GetUserById(org.Members[i].UserID)
			// Update member values
			org.Members[i].Name = user.Name
			org.Members[i].Email = user.Email

		}
	}

	return organizations, nil
}
func GetOrganizationById(OrgID string) (*models.Organization, error) {
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
	objID, err := primitive.ObjectIDFromHex(OrgID)
	if err != nil {
		return nil, err
	}

	// Define a filter to find an organization by ID
	filter := bson.M{"_id": objID}

	// Execute MongoDB operation to find an organization by ID
	var organization models.Organization
	err = collection.FindOne(context.Background(), filter).Decode(&organization)
	if err != nil {
		return nil, err
	}

	// Loop through each member_id in the organization
	for i := range organization.Members {
		user, _ := users.GetUserById(organization.Members[i].UserID)
		// Update member values
		organization.Members[i].Name = user.Name
		organization.Members[i].Email = user.Email
	}

	return &organization, nil
}
func InsertOrganization(NewOrg models.OrganizationBind, UserID string, AccessLevel string) (string, error) {
	// Initialize the database
	if err := InitDB(DBName, CollectionName); err != nil {
		return "", err
	}

	// make sure to close db after this function completes
	defer func() {
		if err := CloseDB(); err != nil {
			fmt.Println("Error closing database:", err)
		}
	}()

	// Execute MongoDB operation to insert a new organization
	// add empty array to member
	NewOrg.Members = []models.OrgMember{}
	result, err := collection.InsertOne(context.Background(), NewOrg)
	if err != nil {
		return "", err
	}
	// Extract the ID from the result and convert it to a string
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	OrgID := objectID.Hex()
	// add creator as fullaccess to organziation member
	err = InsertMemberIntoOrganization(OrgID, UserID, AccessLevel)
	if err != nil {
		return "", err
	}
	return OrgID, nil
}
func InsertMemberIntoOrganization(orgID string, UserID string, AccessLevel string) error {
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

	// cast string to mongo obj id type
	objID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return err
	}
	// Define a filter to find the organization by ID
	filter := bson.M{"_id": objID}

	// Define an update to add the new member into the Members array
	// $addToSet operator adds only if they are not already present
	update := bson.M{
		"$addToSet": bson.M{"members": bson.M{"userid": UserID, "accesslevel": AccessLevel}},
	}

	// Execute MongoDB operation to update the organization
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func InviteMemberToOrganization(orgID string, Email string) error {
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

	user, _ := users.GetUserByEmail(Email)
	err := InsertMemberIntoOrganization(orgID, user.UserID, "readonly")
	if err != nil {
		return err
	}
	return nil
}
func UpdateOrganization(OrgID string, updatedOrg models.OrganizationOnly) error {
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

	objID, err := primitive.ObjectIDFromHex(OrgID)
	if err != nil {
		return err
	}

	// Execute MongoDB operation to update the user
	update := bson.M{"$set": bson.M{"name": updatedOrg.Name, "description": updatedOrg.Description}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOrganizationById(OrgID string) error {
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

	objID, err := primitive.ObjectIDFromHex(OrgID)
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
