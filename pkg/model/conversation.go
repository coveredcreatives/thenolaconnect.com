package model

import (
	"time"

	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
)

type Conversation struct {
	AccountSid          string    `json:"account_sid" gorm:"column:account_sid"`
	ChatServiceSid      string    `json:"chat_service_sid" gorm:"column:chat_service_sid"`
	MessagingServiceSid string    `json:"messaging_service_sid" gorm:"column:messaging_service_sid"`
	Sid                 string    `json:"sid" gorm:"column:sid;primaryKey"`
	FriendlyName        string    `json:"friendly_name" gorm:"column:friendly_name"`
	UniqueName          string    `json:"unique_name" gorm:"column:unique_name"`
	Attributes          string    `json:"attributes" gorm:"column:attributes"`
	State               string    `json:"state" gorm:"column:state"`
	Url                 string    `json:"url" gorm:"column:url"`
	CreatedAt           time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt           time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (c *Conversation) TableName() string {
	return "twilio.conversation"
}

func ConversationToLocalSchema(in conversations_openapi.ConversationsV1ServiceConversation) Conversation {
	out := Conversation{}
	if in.AccountSid != nil {
		out.AccountSid = *in.AccountSid
	}
	if in.ChatServiceSid != nil {
		out.ChatServiceSid = *in.ChatServiceSid
	}
	if in.MessagingServiceSid != nil {
		out.MessagingServiceSid = *in.MessagingServiceSid
	}
	if in.Sid != nil {
		out.Sid = *in.Sid
	}
	if in.FriendlyName != nil {
		out.FriendlyName = *in.FriendlyName
	}
	if in.UniqueName != nil {
		out.UniqueName = *in.UniqueName
	}
	if in.Attributes != nil {
		out.Attributes = *in.Attributes
	}
	if in.State != nil {
		out.State = *in.State
	}
	if in.Url != nil {
		out.Url = *in.Url
	}
	return out
}
