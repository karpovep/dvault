package services

import "dvault/models"

type UserService interface {
	CreateUser(user *models.User) error
	GetUser(string) (*models.User, error)
}
