package repo

import (
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
)

type BusinessRepo struct {
	db *gorm.DB
}

func NewBusinessRepo(db *gorm.DB) *BusinessRepo {
	return &BusinessRepo{db: db}
}

func (br *BusinessRepo) CreateBusiness(user_id int64, business model.Business) error {
	err := br.db.Create(&business).Error
	if err != nil {
		return err
	}

	m := model.UserToBusiness{
		UserID:     user_id,
		BusinessID: business.ID,
	}

	err = br.db.Create(&m).Error
	if err != nil {
		return err
	}

	return nil
}

func (br *BusinessRepo) GetAllBusinesses() ([]model.Business, error) {
	var businesses []model.Business
	return businesses, br.db.Find(&businesses).Error
}

func (br *BusinessRepo) GetBusinessByID(id int64) (model.Business, error) {
	var business model.Business
	return business, br.db.Where("id = ?", id).First(&business).Error
}

func (br *BusinessRepo) DeleteBusiness(id int64) error {
	return br.db.Where("id = ?", id).Delete(&model.Business{}).Error
}

func (br *BusinessRepo) UpdateBusiness(id int64, business model.Business) error {
	return br.db.Model(&model.Business{}).Where("id = ?", id).Updates(business).Error
}

func (br *BusinessRepo) BusinessExists(id int64) (bool, error) {
	var business model.Business
	return business.ID != 0, br.db.Where("id = ?", id).First(&business).Error
}

func (br *BusinessRepo) GetBusinessByINN(inn int64) (model.Business, error) {
	var business model.Business
	return business, br.db.Where("inn = ?", inn).First(&business).Error
}

func (br *BusinessRepo) GetBusinessByOGRN(ogrn int64) (model.Business, error) {
	var business model.Business
	return business, br.db.Where("ogrn = ?", ogrn).First(&business).Error
}

func (br *BusinessRepo) GetUserBusinesses(userID int64) ([]model.Business, error) {
	var businesses []model.Business

	err := br.db.Joins("JOIN user_to_businesses ON businesses.id = user_to_businesses.business_id").
		Where("user_to_businesses.user_id = ?", userID).
		Find(&businesses).Error

	if err != nil {
		return nil, err
	}
	return businesses, nil
}

func (br *BusinessRepo) GetBusinessesUsers(userID int64) ([]model.User, error) {
	var users []model.User

	err := br.db.Joins("JOIN user_to_businesses ON users.id = user_to_businesses.user_id").
		Where("user_to_businesses.business_id = ?", userID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (br *BusinessRepo) AddUserToBusiness(userID int64, businessID int64) error {
	m := model.UserToBusiness{
		UserID:     userID,
		BusinessID: businessID,
	}

	return br.db.Create(&m).Error
}

func (br *BusinessRepo) RemoveUserFromBusiness(userID int64, businessID int64) error {
	return br.db.Where("user_id = ? AND business_id = ?", userID, businessID).Delete(&model.UserToBusiness{}).Error
}
