package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	OrderId                     int            `json:"order_id" gorm:"column:order_id;primaryKey;autoincrement;default:1"`
	ParentOrderId               int            `json:"parent_order_id" gorm:"column:parent_order_id"`
	PDFGenerationWorkerId       int            `json:"pdf_generation_worker_id" gorm:"column:pdf_generation_worker_id"`
	IsPDFGenerated              bool           `json:"is_pdf_generated" gorm:"column:is_pdf_generated"`
	OrderFileStorageURL         string         `json:"order_file_storage_url" gorm:"column:order_file_storage_url"`
	IsViewedByManager           bool           `json:"is_viewed_by_manager" gorm:"column:is_viewed_by_manager"`
	IsAcceptedByManager         bool           `json:"is_accepted_by_manager" gorm:"column:is_accepted_by_manager"`
	IsDeliveredToKitchen        bool           `json:"is_delivered_to_kitchen" gorm:"column:is_delivered_to_kitchen"`
	ConversationSid             string         `json:"conversation_sid" gorm:"column:conversation_sid"`
	FirstConversationMessageSid string         `json:"first_conversation_message_sid" gorm:"column:first_conversation_message_sid"`
	IncomingPhoneNumber         string         `json:"incoming_phone_number" gorm:"column:incoming_phone_number"`
	FormId                      string         `json:"form_id" gorm:"column:form_id"`
	FormResponseId              string         `json:"form_response_id" gorm:"column:form_response_id"`
	MediaSid                    string         `json:"media_sid" gorm:"column:media_sid"`
	CreatedAt                   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                   time.Time      `json:"upated_at" gorm:"coulmn:updated_at"`
	DeletedAt                   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (i *Order) TableName() string {
	return "order_communication.order"
}

func (i *Order) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&i)
}

func (i *Order) Delete(tx *gorm.DB) *gorm.DB {
	if i.OrderId > 0 {
		return tx.Delete(&i)
	}
	return tx
}

func (i *Order) FindOne(tx *gorm.DB) *gorm.DB {
	return tx.Model(&i).Find(&i)
}

func (i *Order) MediaFilename() string {
	if i.FormResponseId != "" {
		return fmt.Sprint(i.FormResponseId, "-", i.OrderId, ".pdf")
	}
	return ""
}
