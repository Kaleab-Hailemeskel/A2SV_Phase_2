package data

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// these litterals should not be accessed anywhere else other than this file
const (
	connectionString = "mongodb://localhost:27017"
	dataBaseName     = "TaskBase"
	collectionName   = "Tasks"
)

// Server Runner
func main() {
	ConnectToMongo()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to close the server(yes/no): ")
		input, _ := reader.ReadString('\n')
		if input == "yes" {
			break
		}
	}
	CloseMongoDB()

}

var mongoClient *mongo.Client

func CloseMongoDB() error {
	return mongoClient.Disconnect(context.TODO())
}
func ConnectToMongo() error {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		return err
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}
	mongoClient = client
	return nil
}

func IsClientConnected() bool {
	return mongoClient != nil
}
func taskAlreadyExists(taskID string) bool {
	_, err := FindByID(taskID)
	return err == nil
}
func FindALL() ([]models.Task, error) {
	collection := mongoClient.Database(dataBaseName).Collection(collectionName)
	resultCursors, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println("Error while fetching All Data from DB", err)
		return nil, err
	}
	tasks := make([]models.Task, 0)
	for resultCursors.Next(context.TODO()) {
		var task models.Task
		err := resultCursors.Decode(&task)
		if err != nil {
			fmt.Println("Error while Decodeing a single data from Cursor", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func FindByID(taskID string) (*models.Task, error) {
	collection := mongoClient.Database(dataBaseName).Collection(collectionName)
	filter := bson.M{"id": taskID}
	var taskResult *models.Task
	err := collection.FindOne(context.TODO(), filter).Decode(&taskResult)
	if err != nil {
		return nil, err
	}
	return taskResult, nil
}
func DeleteOne(taskID string) error {
	collection := mongoClient.Database(dataBaseName).Collection(collectionName)
	filter := bson.M{"id": taskID}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error While Deleteing a task")
		return err
	}
	if deleteResult.DeletedCount == 0 {
		fmt.Println("No Match found")
		return errors.New("task with ID " + taskID + " isn't found")
	}
	return nil
}
func UpdateOne(taskID string, updatedTask models.Task) error {
	if taskID != updatedTask.ID && taskAlreadyExists(updatedTask.ID) {
		return errors.New("task With The new ID already exists, Use Unique ID")
	}
	coll := mongoClient.Database(dataBaseName).Collection(collectionName)
	filter := bson.M{"id": taskID}

	updateFilter := bson.M{"$set": bson.M{
		"id":          updatedTask.ID,
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"due_date":    updatedTask.DueDate,
		"status":      updatedTask.Status}}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, updateFilter)
	if err != nil {
		fmt.Println("cannot update in Database")
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no Data Modified")
	}
	fmt.Println("Updated with ID", updateResult.UpsertedID, "\t", updateResult)

	return nil
}
func InsertOne(t models.Task) error {
	if taskAlreadyExists(t.ID) {
		return errors.New("task with ID < " + t.ID + " > Already exists.")
	}
	collection := mongoClient.Database(dataBaseName).Collection(collectionName)
	inserted, err := collection.InsertOne(context.TODO(), t)
	if err != nil {
		fmt.Println("Single data Insertion Error")
		return err
	}
	fmt.Println("Data Inserted with ID", inserted.InsertedID)
	return nil
}
