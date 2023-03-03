package model

import (
	"time"

	messaging_openapi "github.com/twilio/twilio-go/rest/messaging/v1"
)

type Service struct {
	AccountSid   string     `json:"account_sid" env:"column:account_sid"`
	Sid          string     `json:"sid" env:"column:sid"`
	FriendlyName string     `json:"friendly_name" env:"column:friendly_name"`
	DateCreated  *time.Time `json:"date_created" env:"column:date_created"`
	DateUpdated  *time.Time `json:"date_updated" env:"column:date_updated"`
}

func (s *Service) TableName() string {
	return "twilio.service"
}

func ServiceToLocalSchema(in messaging_openapi.MessagingV1Service) Service {
	out := Service{}
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
	return out
}
