package model

import (
	"fmt"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	models := []any{
		User{},
		VerificationCode{},
		Business{},
		Product{},
		ProductImage{},
		ProductSpecification{},
		ProductReview{},
		CartItem{},
		Order{},
		OrderItem{},
		UserToBusiness{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(&m); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
