package services

import (
	"context"
	"dazer/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteServiceImpl struct {
	noteCollection *mongo.Collection
	Ctx            context.Context
}

func NewNoteService(collection *mongo.Collection, ctx context.Context) NoteService {
	return &NoteServiceImpl{
		noteCollection: collection,
		Ctx:            ctx,
	}
}

func (n *NoteServiceImpl) CreateNote(note *models.Note) error {
	toInsert := &models.Note{
		ID:      primitive.NewObjectID(),
		Title:   note.Title,
		Content: note.Content,
		UserID:  note.UserID,
	}
	_, err := n.noteCollection.InsertOne(n.Ctx, toInsert)
	return err
}

func (n *NoteServiceImpl) GetNote(id string) (*models.Note, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var note = &models.Note{}
	query := bson.D{bson.E{Key: "_id", Value: _id}}
	err = n.noteCollection.FindOne(n.Ctx, query).Decode(note)
	return note, err
}

func (n *NoteServiceImpl) GetAll(userId string) ([]*models.Note, error) {
	var notes []*models.Note
	cursor, err := n.noteCollection.Find(n.Ctx, bson.D{bson.E{Key: "user_id", Value: userId}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(n.Ctx) {
		var note models.Note
		err := cursor.Decode(&note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	if err := cursor.Close(n.Ctx); err != nil {
		return nil, err
	}
	return notes, nil
}

func (n *NoteServiceImpl) UpdateNote(note *models.Note) error {
	filter := bson.D{bson.E{Key: "_id", Value: note.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "title", Value: note.Title},
			bson.E{Key: "content", Value: note.Content}},
		}}
	result, err := n.noteCollection.UpdateOne(n.Ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return errors.New("no docs were found")
	}
	return nil
}

func (n *NoteServiceImpl) DeleteNote(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	query := bson.D{bson.E{Key: "_id", Value: _id}}
	result, err := n.noteCollection.DeleteOne(n.Ctx, query)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("no docs were found")
	}
	return nil
}
