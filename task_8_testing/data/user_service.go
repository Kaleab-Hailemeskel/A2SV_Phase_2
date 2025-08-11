package data

import (
	"context"
	"log"
	"task_8_testing/infrastructure"
	"task_8_testing/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (taskDB *UserDB) CloseDataBase() error {
	return taskDB.Coll.Database().Client().Disconnect(taskDB.Contxt)
}
func (userDB *UserDB) FindUserByID(userID primitive.ObjectID) (*models.UserDTO, error) {

	filter := bson.M{"_id": userID}
	var user models.UserDTO
	err := userDB.Coll.FindOne(userDB.Contxt, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (userDB *UserDB) FindUserByEmail(userEmail string) (*models.UserDTO, error) {

	filter := bson.M{"email": userEmail}
	var user models.UserDTO
	err := userDB.Coll.FindOne(userDB.Contxt, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (userDB *UserDB) StoreUser(user *models.User) (*models.UserDTO, error) {
	insert, err := userDB.Coll.InsertOne(userDB.Contxt, user)
	if err != nil {
		return nil, err
	}
	userDTO := models.ChangeUserModel(user)
	userDTO.ID = insert.InsertedID.(primitive.ObjectID)
	return userDTO, nil
}
func (userDB *UserDB) CheckUserExistance(userEmail string) bool {
	_, err := userDB.FindUserByEmail(userEmail)
	return err == nil
}
func (userDB *UserDB) CheckUserExistanceByID(userID primitive.ObjectID) bool {
	_, err := userDB.FindUserByID(userID)
	return err == nil
}
