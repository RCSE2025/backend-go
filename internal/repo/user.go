package repo

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(user model.User) (model.User, error) {
	return user, r.db.Create(&user).Error
}

func (r *UserRepo) DeleteUser(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *UserRepo) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	return user, r.db.Where("email = ?", email).First(&user).Error
}

func (r *UserRepo) GetUserByID(id int64) (model.User, error) {
	var user model.User
	return user, r.db.Where("id = ?", id).First(&user).Error
}

func (r *UserRepo) GetAllUsers() ([]model.User, error) {
	var users []model.User
	return users, r.db.Find(&users).Error
}

func (r *UserRepo) EmailExists(email string) (bool, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	ok := user.ID != 0

	fmt.Println(ok, err, user)
	return ok, err
}

func (r *UserRepo) UserExists(id int64) (bool, error) {
	var user model.User
	return user.ID != 0, r.db.Where("id = ?", id).First(&user).Error
}

func (r *UserRepo) CreateVerificationCode(code string, expiredAt time.Time, user model.User) (model.VerificationCode, error) {
	verifCode := model.VerificationCode{
		Code:      code,
		UserID:    user.ID,
		ExpiredAt: expiredAt,
	}
	return verifCode, r.db.Create(&verifCode).Error
}

func (r *UserRepo) GetVerificationCode(userID int64, code string) (model.VerificationCode, error) {
	var verifCode model.VerificationCode
	return verifCode, r.db.Where("user_id = ? AND code = ?", userID, code).
		Order("expired_at DESC").
		First(&verifCode).Error
}

func (r *UserRepo) DeleteVerificationCode(code model.VerificationCode) error {
	return r.db.Delete(&code).Error
}

func (r *UserRepo) VerifyEmail(userID int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("is_email_verified", true).Error
}

func (r *UserRepo) SetPassword(userID int64, passwordHash string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error
}

func (r *UserRepo) UpdateUser(userID int64, user model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(user).Error
}
