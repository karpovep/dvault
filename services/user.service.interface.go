package services

import (
	models "dvault/db/entities"
	"dvault/models/dtos"
)

type IUserService interface {
	CreateUser(user *models.User) error
	GetUser(userPubId string) (*models.User, error)
	UpdateUser(userPubId string, dto dtos.UserUpdateRequestDto) (bool, error)
	SearchUsers(search string) ([]*dtos.UserSearchItemResponseDto, error)
}
