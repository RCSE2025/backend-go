package model

import "time"

// User - модель пользователя с тегами JSON и GORM
type User struct {
	BaseModel
	ID                int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name              string       `json:"name" gorm:"size:100;not null"`
	Patronymic        string       `json:"patronymic" gorm:"size:100"`
	Surname           string       `json:"surname" gorm:"size:100;not null"`
	Email             string       `json:"email" gorm:"unique;not null"`
	PasswordHash      string       `json:"-" gorm:"not null"` // "-" исключает поле из JSON
	DateOfBirth       time.Time    `json:"date_of_birth" gorm:"null"`
	IsEmailVerified   bool         `json:"is_email_verified" gorm:"default:false"`
	Role              UserRoleType `json:"role" gorm:"default:user" swaggertype:"primitive,string"`
	IsPasportVerified bool         `json:"is_pasport_verified" gorm:"default:false"`
	INN               *int64       `json:"inn,omitempty" gorm:"unique,null"`
}

type UserRoleType string

const UserRole UserRoleType = "user"
const SelfEmployedRole UserRoleType = "self-employed"
const BusinessRole UserRoleType = "business"
const AdminRole UserRoleType = "admin"
const SupportRole UserRoleType = "support"

type UserToBusiness struct {
	UserID     int64 `json:"user_id" gorm:"not null"`
	BusinessID int64 `json:"business_id" gorm:"not null"`
}

// TableName переопределяет название таблицы (если нужно)
func (User) TableName() string {
	return "users"
}

type UserCreate struct {
	Name              string    `json:"name" `
	Patronymic        string    `json:"patronymic" `
	Surname           string    `json:"surname" `
	Email             string    `json:"email" `
	Password          string    `json:"password" `
	DateOfBirth       time.Time `json:"date_of_birth" `
	IsPasportVerified bool      `json:"is_pasport_verified" `
	INN               *int64    `json:"inn,omitempty" `
}
