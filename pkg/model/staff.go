package model

import (
	"time"

	"gorm.io/gorm"
)

type Staff struct {
	StaffId                    int            `json:"staff_id" gorm:"column:staff_id"`
	Email                      string         `json:"email" gorm:"column:email"`
	PhoneNumber                string         `json:"phone_number" gorm:"column:phone_number"`
	Name                       string         `json:"name" gorm:"column:name"`
	Password                   string         `json:"password" gorm:"column:password"`
	EmailVerificationCode      string         `json:"email_verification_code" gorm:"column:email_verification_code"`
	EmailVerificationCompleted bool           `json:"email_verification_completed" gorm:"column:email_verification_completed"`
	TwoFactorAuthCode          string         `json:"two_factor_auth_code" gorm:"column:two_factor_auth_code"`
	TwoFactorAuthCompleted     bool           `json:"two_factor_auth_completed" gorm:"column:two_factor_auth_completed"`
	IsManager                  bool           `json:"is_manager" gorm:"column:is_manager"`
	CreatedAt                  time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                  time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt                  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (s *Staff) TableName() string {
	return "public.staff"
}
