package repositories

import (
	"context"
	models "dvault/db/entities"
	"dvault/models/dtos"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	result := repo.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return repo.GetByUserPubId(user.UserPubId)
}

func (repo *UserRepository) GetByUserPubId(userPubId string) (*models.User, error) {
	var user models.User
	result := repo.db.First(&user, "user_pub_id = ?", userPubId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) Update(userPubId string, dto dtos.UserUpdateRequestDto) (bool, error) {
	toUpdate := map[string]interface{}{}
	if dto.Username != nil {
		toUpdate["username"] = dto.Username
	}
	if dto.IsPublic != nil {
		toUpdate["is_public"] = dto.IsPublic
	}
	res := repo.db.Model(&models.User{}).
		Where("user_pub_id = ?", userPubId).
		Updates(toUpdate)
	if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}

func (repo *UserRepository) Search(search string) ([]*dtos.UserSearchItemResponseDto, error) {
	var users []*dtos.UserSearchItemResponseDto
	repo.db.Table("users").
		Select("username", "user_pub_id").
		Where("username ILIKE ?", fmt.Sprintf("%%%s%%", search)).
		Limit(5).
		Scan(&users)
	return users, nil
}

func NewUserRepository(ctx context.Context, db *gorm.DB) IUserRepository {
	return &UserRepository{
		ctx,
		db,
	}
}
