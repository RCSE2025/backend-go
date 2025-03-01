package service

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/RCSE2025/backend-go/internal/email"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
	"github.com/RCSE2025/backend-go/internal/utils"
	"github.com/RCSE2025/backend-go/pkg/logger/sl"
	"html/template"
	"log/slog"
	"math/big"
	"time"
)

type UserService struct {
	repo        *repo.UserRepo
	jwtService  JWTService
	mailer      *email.Mailer
	frontendURL string
}

func NewUserService(repo *repo.UserRepo, jwtService JWTService, mailer *email.Mailer) *UserService {
	return &UserService{
		repo:       repo,
		jwtService: jwtService,
		mailer:     mailer,
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
			Role:         model.UserRole,
		},
	)

	code, err := s.CreateVerificationCode(userDB)
	if err != nil {
		return model.User{}, err
	}
	body, err := makeVerificationEmailTemplate(user.Name, code.Code)

	if err != nil {
		return model.User{}, err
	}

	go func() {
		err = s.mailer.SendMail(user.Email, "Подтверждение почты", body)
		if err != nil {
			slog.Error("cannot send verification email", sl.Err(err))
		}
	}()

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

	if !ok {
		return model.Token{}, ErrWrongEmailOrPassword
	}

	token, err := s.GenerateNewToken(user)
	if err != nil {
		return model.Token{}, err
	}
	return token, nil
}

func (s *UserService) GetUserRole(user model.User) string {
	return string(user.Role)
}
func (s *UserService) GenerateNewToken(user model.User) (model.Token, error) {
	accessToken, err := s.jwtService.GenerateToken(user.ID, s.GetUserRole(user))
	if err != nil {
		return model.Token{}, err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, s.GetUserRole(user))
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

const verificationCodeExpirationTime = 10 * time.Minute

func (s *UserService) CreateVerificationCode(user model.User) (model.VerificationCode, error) {

	code, err := generateCode()
	if err != nil {
		return model.VerificationCode{}, err
	}
	return s.repo.CreateVerificationCode(code, time.Now().Add(verificationCodeExpirationTime), user)
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

const verificationEmailTemplate = `
<html>
    <body style="font-family: Arial, sans-serif; color: #333; background-color: #f9f9f9; padding: 20px;">
        <div style="max-width: 600px; margin: auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); text-align: center;">
            <h1 style="color: #4CAF50;">Подтверждение почты</h1>
            <p style="font-size: 16px; line-height: 1.5;">Здравствуйте, {{.Name}}!</p>
            <p style="font-size: 16px; line-height: 1.5;">Ваш код подтверждения:</p>
            <p style="font-size: 24px; font-weight: bold; color: #4CAF50;">{{.Code}}</p>
            <p style="font-size: 14px; color: #666;">Введите этот код в соответствующее поле, чтобы подтвердить свою электронную почту.</p>
        </div>
    </body>
</html>`

func makeVerificationEmailTemplate(name string, code string) (string, error) {
	tmpl, err := template.New("verification_email").Parse(verificationEmailTemplate)
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Name string
		Code string
	}{Name: name, Code: code})

	return body.String(), err
}

var ErrCodeExpired = errors.New("code expired")

func (s *UserService) VerifyEmail(userID int64, code string) error {

	if err := s.UserNotExistsWithErr(userID); err != nil {
		return err
	}

	codeDB, err := s.repo.GetVerificationCode(userID, code)
	if err != nil {
		return err
	}

	if time.Now().After(codeDB.ExpiredAt) {
		return ErrCodeExpired
	}

	if err := s.repo.DeleteVerificationCode(codeDB); err != nil {
		return err
	}

	if err := s.repo.VerifyEmail(userID); err != nil {
		return err
	}
	return nil
}

const resetPasswordEmailTemplate = `
<html>
    <body style="font-family: Arial, sans-serif; color: #333; background-color: #f9f9f9; padding: 20px;">
        <div style="max-width: 600px; margin: auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); text-align: center;">
            <h1 style="color: #4CAF50;">Восстановление пароля</h1>
            <p style="font-size: 16px; line-height: 1.5;">Здравствуйте, {{.Name}}!</p>
            <p style="font-size: 16px; line-height: 1.5;">Вы запросили восстановление пароля. Для установки нового пароля перейдите по ссылке ниже:</p>
            <p>
                <a href="{{.ResetLink}}" style="display: inline-block; padding: 12px 20px; font-size: 16px; color: #fff; background-color: #4CAF50; text-decoration: none; border-radius: 5px;">Сбросить пароль</a>
            </p>
            <p style="font-size: 14px; color: #666;">Если вы не запрашивали восстановление пароля, просто проигнорируйте это письмо.</p>
        </div>
    </body>
</html>
`

func makeResetPasswordEmailTemplate(name string, resetLink string) (string, error) {
	tmpl, err := template.New("reset_password_email").Parse(resetPasswordEmailTemplate)
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Name      string
		ResetLink string
	}{Name: name, ResetLink: resetLink})

	return body.String(), err
}

func (s *UserService) SendResetPasswordEmail(email string) error {
	if err := s.EmailNotExistsWithErr(email); err != nil {
		return err
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	token, err := s.jwtService.GenerateRefreshPasswordToken(user.ID)

	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)

	body, err := makeResetPasswordEmailTemplate(user.Name, resetLink)
	if err != nil {
		return err
	}

	go func() {
		if err := s.mailer.SendMail(user.Email, "Восстановление пароля", body); err != nil {
			slog.Error("cannot send verification email", sl.Err(err))

		}
	}()

	return nil
}

func (s *UserService) UpdateUser(userID int64, user model.User) error {
	return s.repo.UpdateUser(userID, user)
}

func (s *UserService) GenerateRefreshPasswordToken(user model.User) (string, error) {
	token, err := s.jwtService.GenerateRefreshPasswordToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) RefreshPassword(token string, password string) error {
	claims, err := s.jwtService.ValidateRefreshPasswordToken(token)
	if err != nil {
		return err
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.repo.SetPassword(claims.UserID, passwordHash)
	if err != nil {
		return err
	}
	return nil
}
