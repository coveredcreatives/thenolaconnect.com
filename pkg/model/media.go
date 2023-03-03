package model

import (
	"time"
)

type Media struct {
	Sid               string    `json:"sid" gorm:"column:sid"`
	ServiceSid        string    `json:"service_sid" gorm:"column:service_sid"`
	DateCreated       time.Time `json:"date_created" gorm:"column:date_created"`
	DateUploadUpdated time.Time `json:"date_upload_updated" gorm:"column:date_upload_updated"`
	DateUpdated       time.Time `json:"date_updated" gorm:"column:date_updated"`
	Size              int       `json:"size" gorm:"column:size"`
	ContentType       string    `json:"content_type" gorm:"column:content_type"`
	Filename          string    `json:"filename" gorm:"column:filename"`
	Author            string    `json:"author" gorm:"column:author"`
	Category          string    `json:"category" gorm:"column:category"`
	MessageSid        string    `json:"message_sid" gorm:"column:message_sid"`
	ChannelSid        string    `json:"channel_sid" gorm:"column:channel_sid"`
	Url               string    `json:"url" gorm:"column:url"`
}

func (f *Media) TableName() string {
	return "twilio.media"
}
