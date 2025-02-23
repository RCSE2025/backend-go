package model

import (
	"time"
)

// User - модель пользователя с тегами JSON и GORM
type User struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string    `json:"name" gorm:"size:100;not null"`
	Patronymic      string    `json:"patronymic" gorm:"size:100"`
	Surname         string    `json:"surname" gorm:"size:100;not null"`
	Email           string    `json:"email" gorm:"unique;not null"`
	PasswordHash    string    `json:"-" gorm:"not null"` // "-" исключает поле из JSON
	DateOfBirth     string    `json:"date_of_birth" gorm:"type:date"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName переопределяет название таблицы (если нужно)
func (User) TableName() string {
	return "users"
}
