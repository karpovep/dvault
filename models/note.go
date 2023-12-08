package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Note struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Title   string             `json:"title" bson:"title"`
	Content interface{}        `json:"content" bson:"content"`
	UserID  string             `json:"userId" bson:"user_id"`
}
