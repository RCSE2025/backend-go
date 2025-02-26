package service

import (
	"errors"
	"fmt"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
	"github.com/RCSE2025/backend-go/internal/utils"
)

type UserService struct {
	repo       *repo.UserRepo
	jwtService JWTService
}

func NewUserService(repo *repo.UserRepo, jwtService JWTService) *UserService {
	return &UserService{
		repo:       repo,
		jwtService: jwtService,
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

var ErrWrongEmailOrPassword = errors.New("wrong email or password")

func (s *UserService) GetToken(email, password string) (model.Token, error) {
	if err := s.EmailNotExistsWithErr(email); err != nil {
		return model.Token{}, ErrWrongEmailOrPassword
	}

	fmt.Println("email exists")

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return model.Token{}, err
	}

	ok, _ := utils.CheckPassword(user.PasswordHash, []byte(password))

	//if err != nil {
	//	return model.Token{}, err
	//}

	if !ok {
		return model.Token{}, ErrWrongEmailOrPassword
	}

	token, err := s.GenerateNewToken(user)
	if err != nil {
		return model.Token{}, err
	}
	return token, nil
}

func (s *UserService) GenerateNewToken(user model.User) (model.Token, error) {
	accessToken, err := s.jwtService.GenerateToken(user.ID, "user")
	if err != nil {
		return model.Token{}, err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, "user")
	if err != nil {
		return model.Token{}, err
	}
	token := model.Token{
		RefreshToken: accessToken,
		AccessToken:  refreshToken,
		TokenType:    "bearer",
	}
	return token, nil
}

func (s *UserService) RefreshToken(refreshToken string) (model.Token, error) {
	token, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil || !token.Valid {
		return model.Token{}, errors.New("invalid token")
	}

	userId, err := s.jwtService.GetUserIDByToken(refreshToken)

	if err != nil {
		return model.Token{}, err
	}

	if ok, _ := s.UserExists(userId); !ok {
		return model.Token{}, ErrUserNotFound
	}

	user, err := s.repo.GetUserByID(userId)

	if err != nil {
		return model.Token{}, err
	}

	if err != nil {
		return model.Token{}, err
	}

	newToken, err := s.GenerateNewToken(user)
	if err != nil {
		return model.Token{}, err
	}
	return newToken, nil
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

var ErrEmailNotExists = errors.New("email not exists")

func (s *UserService) EmailNotExistsWithErr(email string) error {
	if exists, _ := s.repo.EmailExists(email); !exists {
		return ErrEmailNotExists
	} else {
		return nil
	}
}

func (s *UserService) EmailExistsWithErr(email string) error {
	if exists, _ := s.repo.EmailExists(email); exists {
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
