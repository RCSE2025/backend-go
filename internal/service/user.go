package service

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
	"github.com/RCSE2025/backend-go/internal/utils"
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user model.UserCreate) (model.User, error) {
	if err := s.EmailExistsWithErr(user.Email); err != nil {
		return model.User{}, err
	}

	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		return model.User{}, err

	}

	userDB, err := s.repo.CreateUser(
		model.User{
			Name:         user.Name,
			Patronymic:   user.Patronymic,
			Surname:      user.Surname,
			Email:        user.Email,
			PasswordHash: passwordHash,
			DateOfBirth:  user.DateOfBirth,
		},
	)

	if err != nil {
		return model.User{}, err
	}

	return userDB, nil
}

func (s *UserService) GetUserByID(id int64) (model.User, error) {
	if err := s.UserNotExistsWithErr(id); err != nil {
		return model.User{}, err
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (model.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(id int64) error {
	if err := s.UserNotExistsWithErr(id); err != nil {
		return err
	}
	return s.repo.DeleteUser(id)
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UserExists(id int64) (bool, error) {
	return s.repo.UserExists(id)
}

func (s *UserService) EmailExists(email string) (bool, error) {
	return s.repo.EmailExists(email)
}

func (s *UserService) EmailExistsWithErr(email string) error {
	if exists, err := s.repo.EmailExists(email); exists {
		fmt.Println("exists", exists, err)
		return ErrEmailExists
	} else {
		return nil
	}
}

func (s *UserService) UserNotExistsWithErr(id int64) error {

	if exists, err := s.repo.UserExists(id); !exists {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}
	return nil
}
