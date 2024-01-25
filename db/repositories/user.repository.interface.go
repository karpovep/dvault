package repositories

import (
	models "dvault/db/entities"
	"dvault/models/dtos"
)

type IUserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetByUserPubId(userPubId string) (*models.User, error)
	Update(userPubId string, dto dtos.UserUpdateRequestDto) (bool, error)
	Search(search string) ([]*dtos.UserSearchItemResponseDto, error)
}
