package data

import (
	"context"
	"errors"
	"log"
	"task_7_clean_architecture/infrastructure"
	"task_7_clean_architecture/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskDB struct {
	Coll   mongo.Collection
	Contxt context.Context
}

func NewTaskDataBaseService() models.ITaskDataBase {
	log.Println(infrastructure.TASK_DB, infrastructure.TASK_COLLECTION_NAME)

	taskDataBaseName := infrastructure.TASK_DB
	taskCollectionName := infrastructure.TASK_COLLECTION_NAME
	connectionString := infrastructure.CONNECTION_STRING

	collection := InitDataBase(connectionString, taskDataBaseName, taskCollectionName)
	return &TaskDB{
		Coll:   *collection,
		Contxt: context.TODO(),
	}

}
func InitDataBase(connectionString, dbName, collectionName string) *mongo.Collection {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(connectionString, " <> ", err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbName).Collection(collectionName)
	return coll
}

func (taskDB *TaskDB) CloseDataBase() error {
	return taskDB.Coll.Database().Client().Disconnect(taskDB.Contxt)
}
func (taskDB *TaskDB) CheckTaskExistance(taskID string) bool {
	_, err := taskDB.FindByID(taskID)
	return err == nil
}
func (taskDB *TaskDB) FindAllTasks(userEmail string) (*[]models.Task, error) {
	var filter primitive.M
	if userEmail != "" {
		filter = bson.M{"owneremail": userEmail}
	} else {
		filter = bson.M{}
	}
	resultCursors, err := taskDB.Coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer resultCursors.Close(context.TODO())

	tasks := make([]models.Task, 0)
	for resultCursors.Next(context.TODO()) {
		var task models.Task
		err := resultCursors.Decode(&task)
		if err != nil {
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil
}
func (taskDB *TaskDB) FindByID(taskID string) (*models.Task, error) {
	filter := bson.M{"id": taskID}
	var taskResult *models.Task
	err := taskDB.Coll.FindOne(taskDB.Contxt, filter).Decode(&taskResult)
	if err != nil {
		return nil, err
	}
	return taskResult, nil
}
func (taskDB *TaskDB) DeleteOne(taskID string) error {
	filter := bson.M{"id": taskID}
	deleteResult, err := taskDB.Coll.DeleteOne(context.TODO(), filter)
	if err != nil {

		return err
	}
	if deleteResult.DeletedCount == 0 {

		return errors.New("task with ID " + taskID + " isn't found")
	}
	return nil
}
func (taskDB *TaskDB) UpdateOne(taskID string, updatedTask models.Task) error {
	if taskID != updatedTask.ID && taskDB.CheckTaskExistance(updatedTask.ID) {
		return errors.New("task With The new ID already exists, Use Unique ID")
	}
	filter := bson.M{"id": taskID}

	updateFilter := bson.M{"$set": bson.M{
		"id":          updatedTask.ID,
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"due_date":    updatedTask.DueDate,
		"status":      updatedTask.Status}}
	updateResult, err := taskDB.Coll.UpdateOne(taskDB.Contxt, filter, updateFilter)
	if err != nil {

		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no Data Modified")
	}
	return nil
}
func (taskDB *TaskDB) InsertOne(t models.Task) error {
	if taskDB.CheckTaskExistance(t.ID) {
		return errors.New("task with ID < " + t.ID + " > Already exists.")
	}
	_, err := taskDB.Coll.InsertOne(taskDB.Contxt, t)
	if err != nil {
		return err
	}
	return nil
}
