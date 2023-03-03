package model

import (
	"encoding/json"
	"time"

	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
	"gorm.io/gorm"
)

type ConversationService struct {
	AccountSid   string         `json:"account_sid" env:"column:account_sid"`
	Sid          string         `json:"sid" env:"column:sid"`
	FriendlyName string         `json:"friendly_name" env:"column:friendly_name"`
	DateCreated  *time.Time     `json:"date_created" env:"column:date_created"`
	DateUpdated  *time.Time     `json:"date_updated" env:"column:date_updated"`
	Url          string         `json:"url" env:"column:url"`
	Links        string         `json:"links" env:"column:links"`
	CreatedAt    time.Time      `json:"created_at" env:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" env:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" env:"column:deleted_at"`
}

func (s *ConversationService) TableName() string {
	return "twilio.conversation_service"
}

func ConversationServiceToLocalSchema(in conversations_openapi.ConversationsV1Service) ConversationService {
	out := ConversationService{}
	if in.AccountSid != nil {
		out.AccountSid = *in.AccountSid
	}
	if in.Sid != nil {
		out.Sid = *in.Sid
	}
	if in.FriendlyName != nil {
		out.FriendlyName = *in.FriendlyName
	}
	if in.DateCreated != nil {
		out.DateCreated = in.DateCreated
	}
	if in.DateUpdated != nil {
		out.DateUpdated = in.DateUpdated
	}
	if in.Url != nil {
		out.Url = *in.Url
	}

	if in.Links != nil {
		bytes, err := json.Marshal(in.Links)
		if err == nil {
			out.Links = string(bytes)
		}
	}
	return out
}
