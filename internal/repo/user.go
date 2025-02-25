package repo

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
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
	ok := user.ID != 0
	err := r.db.Where("email = ?", email).First(&user).Error
	fmt.Println(ok, err, user)
	return ok, err
}

func (r *UserRepo) UserExists(id int64) (bool, error) {
	var user model.User
	return user.ID != 0, r.db.Where("id = ?", id).First(&user).Error
}
