package services

import (
	"context"
	"dvault/app"
	"dvault/constants"
	models "dvault/db/entities"
	"dvault/db/repositories"
	"dvault/models/dtos"
)

type UserService struct {
	Ctx            context.Context
	UserRepository repositories.IUserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
	_, err := s.UserRepository.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUser(userPubId string) (*models.User, error) {
	return s.UserRepository.GetByUserPubId(userPubId)
}

func (s *UserService) UpdateUser(userPubId string, dto dtos.UserUpdateRequestDto) (bool, error) {
	return s.UserRepository.Update(userPubId, dto)
}

func (s *UserService) SearchUsers(search string) ([]*dtos.UserSearchItemResponseDto, error) {
	return s.UserRepository.Search(search)
}

func NewUserService(appContext app.IAppContext) IUserService {
	ctx := appContext.Get(constants.Ctx).(context.Context)
	userRepository := appContext.Get(constants.UserRepository).(repositories.IUserRepository)
	return &UserService{
		ctx,
		userRepository,
	}
}
