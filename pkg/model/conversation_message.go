package model

import (
	"encoding/json"
	"time"

	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"

	"gorm.io/gorm"
)

type ConversationMessage struct {
	AccountSid      string         `json:"account_sid" gorm:"column:account_sid"`
	ConversationSid string         `json:"conversation_sid" gorm:"column:conversation_sid"`
	Sid             string         `json:"sid" gorm:"column:sid"`
	Index           int            `json:"index" gorm:"column:index"`
	Author          string         `json:"author" gorm:"column:author"`
	Body            string         `json:"body" gorm:"column:body"`
	Media           string         `json:"media" gorm:"column:media"`
	Attributes      string         `json:"attributes" gorm:"column:attributes"`
	ParticipantSid  string         `json:"participant_sid" gorm:"column:participant_sid"`
	Url             string         `json:"url" gorm:"column:url"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (i *ConversationMessage) TableName() string {
	return "twilio.conversation_message"
}

func ConversationMessageToSchema(in conversations_openapi.ConversationsV1ServiceConversationMessage) ConversationMessage {
	out := ConversationMessage{}
	if in.AccountSid != nil {
		out.AccountSid = *in.AccountSid
	}
	if in.ConversationSid != nil {
		out.ConversationSid = *in.ConversationSid
	}
	if in.Sid != nil {
		out.Sid = *in.Sid
	}
	if in.Index != nil {
		out.Index = *in.Index
	}
	if in.Author != nil {
		out.Author = *in.Author
	}
	if in.Body != nil {
		out.Body = *in.Body
	}
	if in.Media != nil {
		marshalled, _ := json.Marshal(in.Media)
		out.Media = string(marshalled)
	}
	if in.Attributes != nil {
		marshalled, _ := json.Marshal(in.Attributes)
		out.Attributes = string(marshalled)
	}
	if in.ParticipantSid != nil {
		out.ParticipantSid = *in.ParticipantSid
	}
	if in.Url != nil {
		out.Url = *in.Url
	}
	return out
}
