package model

import (
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	models := []any{User{}, VerificationCode{}}

	for _, m := range models {
		if err := db.AutoMigrate(&m); err != nil {
			return err
		}
	}

	return nil
}
