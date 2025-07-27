package data

import (
	"context"
	"log"
	"task_7_clean_architecture/infrastructure"
	"task_7_clean_architecture/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDB struct {
	Coll   *mongo.Collection
	Contxt context.Context
}

func NewUserDataBase() models.IUserDataBase {
	userDataBase := infrastructure.USER_DB
	userCollection := infrastructure.USER_COLLECTION_NAME
	connectionString := infrastructure.CONNECTION_STRING
	log.Println("while InitInUserDataBase: ", userDataBase, connectionString, userCollection)
	return &UserDB{
		Coll:   InitDataBase(connectionString, userDataBase, userCollection),
		Contxt: context.TODO(),
	}
}
func (taskDB *UserDB) CloseDataBase() error { //! a function with The same parameter exists in the task_service.go
	return taskDB.Coll.Database().Client().Disconnect(taskDB.Contxt)
}
func (userDB *UserDB) FindUserByEmail(userEmail string) (*models.User, error) {

	filter := bson.M{"email": userEmail}
	var user models.User
	err := userDB.Coll.FindOne(userDB.Contxt, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (userDB *UserDB) StoreUser(user *models.User) error {
	_, err := userDB.Coll.InsertOne(userDB.Contxt, user)
	if err != nil {
		return err
	}
	return nil
}
func (userDB *UserDB) CheckUserExistance(userEmail string) bool {
	_, err := userDB.FindUserByEmail(userEmail)
	return err == nil
}
