package model

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime,nullable"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime,nullable"`
}
