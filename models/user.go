package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID string             `json:"userId" bson:"user_id"`
}
