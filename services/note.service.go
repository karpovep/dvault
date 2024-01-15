package services

import "dvault/models"

type NoteService interface {
	CreateNote(note *models.Note) (*models.Note, error)
	GetNote(string) (*models.Note, error)
	GetAll(string) ([]*models.Note, error)
	UpdateNote(*models.Note) (*models.Note, error)
	DeleteNote(string) error
}
