package services

import (
	"context"
	"dazer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	Ctx            context.Context
}

func NewUserService(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		userCollection: collection,
		Ctx:            ctx,
	}
}

func (service *UserServiceImpl) CreateUser(user *models.User) error {
	toInsert := &models.User{
		ID:     primitive.NewObjectID(),
		UserID: user.UserID,
	}
	_, err := service.userCollection.InsertOne(service.Ctx, toInsert)
	return err
}

func (service *UserServiceImpl) GetUser(userId string) (*models.User, error) {
	var user = &models.User{}
	query := bson.D{bson.E{Key: "user_id", Value: userId}}
	err := service.userCollection.FindOne(service.Ctx, query).Decode(user)
	return user, err
}
