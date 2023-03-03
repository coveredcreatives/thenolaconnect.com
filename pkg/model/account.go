package model

import (
	"time"

	api_openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"gorm.io/gorm"
)

type Account struct {
	Sid             string         `json:"sid" gorm:"column:sid"`
	OwnerAccountSid string         `json:"owner_account_sid" gorm:"column:owner_account_sid"`
	AuthToken       string         `json:"auth_token" gorm:"column:auth_token"`
	FriendlyName    string         `json:"friendly_name" gorm:"column:friendly_name"`
	Status          string         `json:"status" gorm:"column:status"`
	Type            string         `json:"type" gorm:"column:type"`
	Uri             string         `json:"uri" gorm:"column:uri"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (a *Account) TableName() string {
	return "twilio.account"
}

func AccountToLocalSchema(in api_openapi.ApiV2010Account) Account {
	out := Account{}
	if in.Sid != nil {
		out.Sid = *in.Sid
	}
	if in.OwnerAccountSid != nil {
		out.OwnerAccountSid = *in.OwnerAccountSid
	}
	if in.AuthToken != nil {
		out.AuthToken = *in.AuthToken
	}
	if in.FriendlyName != nil {
		out.FriendlyName = *in.FriendlyName
	}
	if in.Status != nil {
		out.Status = *in.Status
	}
	if in.Type != nil {
		out.Type = *in.Type
	}
	if in.Uri != nil {
		out.Uri = *in.Uri
	}
	return out
}
