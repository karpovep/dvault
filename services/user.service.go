package services

import "dazer/models"

type UserService interface {
	CreateUser(user *models.User) error
	GetUser(string) (*models.User, error)
}
