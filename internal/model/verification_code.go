package model

import (
	"time"
)

// VerificationCode - .
type VerificationCode struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	Code      string    `json:"code" gorm:"size:10;not null"`
	SentAt    time.Time `json:"sent_at" gorm:"autoCreateTime"`
	ExpiredAt time.Time `json:"expired_at"`
}

// TableName - .
func (VerificationCode) TableName() string {
	return "verification_codes"
}
