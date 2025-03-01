package repo

import (
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(item model.Product) (model.Product, error) {
	return item, r.db.Create(&item).Error
}

func (r *ProductRepo) FindByID(id int64) (model.Product, error) {
	var product model.Product
	return product, r.db.Where("id = ?", id).First(&product).Error
}

func (r *ProductRepo) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Product{}).Error
}
