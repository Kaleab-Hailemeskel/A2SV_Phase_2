package data

import (
	"context"
	"errors"
	"log"
	"task_8_testing/infrastructure"
	"task_8_testing/models"

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
func (taskDB *TaskDB) CheckTaskExistance(taskID primitive.ObjectID) bool {
	_, err := taskDB.FindByID(taskID)
	return err == nil
}
func (taskDB *TaskDB) FindAllTasks(userEmail string) ([]*models.TaskDTO, error) {
	log.Println("🦌 Deep into the DB")
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

	tasks := make([]*models.TaskDTO, 0)
	for resultCursors.Next(context.TODO()) {
		var task models.TaskDTO
		err := resultCursors.Decode(&task)
		if err != nil {
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
func (taskDB *TaskDB) FindByID(taskID primitive.ObjectID) (*models.TaskDTO, error) {
	filter := bson.M{"id": taskID}
	var taskResult *models.TaskDTO
	err := taskDB.Coll.FindOne(taskDB.Contxt, filter).Decode(&taskResult)
	if err != nil {
		return nil, err
	}
	return taskResult, nil
}
func (taskDB *TaskDB) DeleteOne(taskID primitive.ObjectID) error {
	filter := bson.M{"id": taskID}
	deleteResult, err := taskDB.Coll.DeleteOne(context.TODO(), filter)
	if err != nil {

		return err
	}
	if deleteResult.DeletedCount == 0 {

		return errors.New("task with ID " + taskID.Hex() + " isn't found")
	}
	return nil
}
func (taskDB *TaskDB) UpdateOne(taskID primitive.ObjectID, updatedTask *models.TaskDTO) (*models.TaskDTO, error) {
	if taskID != updatedTask.ID && taskDB.CheckTaskExistance(updatedTask.ID) {
		return nil, errors.New("task With The new ID already exists, Use Unique ID")
	}
	filter := bson.M{"_id": taskID}

	updateFilter := bson.M{"$set": bson.M{
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"due_date":    updatedTask.DueDate,
		"status":      updatedTask.Status}}
	updateResult, err := taskDB.Coll.UpdateOne(taskDB.Contxt, filter, updateFilter)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no Data Modified")
	}
	return updatedTask, nil
}
func (taskDB *TaskDB) InsertOne(t *models.TaskDTO) (*models.TaskDTO, error) {
	if taskDB.CheckTaskExistance(t.ID) {
		return nil, errors.New("task with ID < " + t.ID.Hex() + " > Already exists.")
	}
	_, err := taskDB.Coll.InsertOne(taskDB.Contxt, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
