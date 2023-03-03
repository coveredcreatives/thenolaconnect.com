package model

import (
	"encoding/json"
	"fmt"
	"time"

	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"

	"gorm.io/gorm"
)

type ConversationParticipant struct {
	AccountSid                   string         `json:"account_sid" gorm:"column:account_sid"`
	ConversationSid              string         `json:"conversation_sid" gorm:"column:conversation_sid"`
	Sid                          string         `json:"sid" gorm:"column:sid"`
	Identity                     string         `json:"identity" gorm:"column:identity"`
	Attributes                   string         `json:"attributes" gorm:"column:attributes"`
	MessagingBinding             string         `json:"messaging_binding" gorm:"column:messaging_binding"`
	MessagingBindingAddress      string         `json:"messaging_binding_address" gorm:"column:messaging_binding_address"`
	MessagingBindingProxyAddress string         `json:"messaging_binding_proxy_address" gorm:"column:messaging_binding_proxy_address"`
	RoleSid                      string         `json:"role_sid" gorm:"column:role_sid"`
	DateCreated                  time.Time      `json:"date_created" gorm:"column:date_created"`
	DateUpdated                  time.Time      `json:"date_updated" gorm:"column:date_updated"`
	Url                          string         `json:"url" gorm:"column:url"`
	LastReadMessageIndex         int            `json:"last_read_message_index" gorm:"column:last_read_message_index"`
	LastReadTimestamp            string         `json:"last_read_timestamp"`
	CreatedAt                    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt                    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (c *ConversationParticipant) TableName() string {
	return "twilio.conversation_participant"
}

func ConversationParticipantToSchema(in conversations_openapi.ConversationsV1ServiceConversationParticipant) ConversationParticipant {
	out := ConversationParticipant{}
	if in.AccountSid != nil {
		out.AccountSid = *in.AccountSid
	}
	if in.ConversationSid != nil {
		out.ConversationSid = *in.ConversationSid
	}
	if in.Sid != nil {
		out.Sid = *in.Sid
	}
	if in.Identity != nil {
		out.Identity = *in.Identity
	}
	if in.Attributes != nil {
		out.Attributes = *in.Attributes
	}
	if in.Url != nil {
		out.Url = *in.Url
	}
	if in.RoleSid != nil {
		out.RoleSid = *in.RoleSid
	}
	if in.MessagingBinding != nil {
		message, _ := json.Marshal(in.MessagingBinding)
		out.MessagingBinding = string(message)
		messaging_binding := map[string]string{"address": "", "proxy_address": ""}
		_ = json.Unmarshal(message, &messaging_binding)
		out.MessagingBindingAddress = messaging_binding["address"]
		out.MessagingBindingProxyAddress = messaging_binding["proxy_address"]
	}
	bits, _ := json.Marshal(out)
	fmt.Println(string(bits))
	return out
}
