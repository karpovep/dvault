package services

import "dazer/models"

type NoteService interface {
	CreateNote(note *models.Note) error
	GetNote(string) (*models.Note, error)
	GetAll(string) ([]*models.Note, error)
	UpdateNote(*models.Note) error
	DeleteNote(string) error
}
