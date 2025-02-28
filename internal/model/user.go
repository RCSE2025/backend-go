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
	DateOfBirth     DateOnly  `json:"date_of_birth" gorm:"type:DATE"  example:"2006-01-02" swaggertype:"primitive,string" format:"date"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"default:false"`
	IsAdmin         bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

//func (u *User) UnmarshalJSON(data []byte) error {
//	type Alias User
//	aux := &struct {
//		DateOfBirth *string `json:"date_of_birth"`
//		*Alias
//	}{
//		Alias: (*Alias)(u),
//	}
//
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//
//	if aux.DateOfBirth != nil && *aux.DateOfBirth != "" {
//		t, err := time.Parse(dateFormat, *aux.DateOfBirth)
//		if err != nil {
//			return err
//		}
//		u.DateOfBirth = &t
//	} else {
//		u.DateOfBirth = nil // Устанавливаем nil, если дата пустая
//	}
//
//	return nil
//}
//
//// Кастомный MarshalJSON для возврата даты в формате YYYY-MM-DD
//func (u User) MarshalJSON() ([]byte, error) {
//	type Alias User
//	aux := &struct {
//		DateOfBirth *string `json:"date_of_birth,omitempty"`
//		*Alias
//	}{
//		Alias: (*Alias)(&u),
//	}
//
//	if u.DateOfBirth != nil {
//		formatted := u.DateOfBirth.Format(dateFormat)
//		aux.DateOfBirth = &formatted
//	}
//
//	return json.Marshal(aux)
//}

//const dateFormat = "2006-01-02"
//
//type Date time.Time
//
//func (d *Date) UnmarshalJSON(data []byte) error {
//	// Убираем кавычки из строки
//	str := string(data)
//	str = str[1 : len(str)-1]
//
//	// Если дата пустая, устанавливаем nil
//	if str == "" || str == "null" {
//		*d = Date(time.Time{})
//		return nil
//	}
//
//	// Парсим дату
//	t, err := time.Parse(dateFormat, str)
//	if err != nil {
//		return err
//	}
//
//	*d = Date(t)
//	return nil
//}
//
//func (d Date) MarshalJSON() ([]byte, error) {
//	if time.Time(d).IsZero() {
//		return []byte(`null`), nil
//	}
//	return json.Marshal(time.Time(d).Format(dateFormat))
//}
//
//// Метод для приведения к типу time.Time
//func (d Date) ToTime() time.Time {
//	return time.Time(d)
//}

// TableName переопределяет название таблицы (если нужно)
func (User) TableName() string {
	return "users"
}

type UserCreate struct {
	Name        string   `json:"name" `
	Patronymic  string   `json:"patronymic" `
	Surname     string   `json:"surname" `
	Email       string   `json:"email" `
	Password    string   `json:"password" `
	DateOfBirth DateOnly `json:"date_of_birth" example:"2006-01-02" swaggertype:"primitive,string" format:"date"`
}
