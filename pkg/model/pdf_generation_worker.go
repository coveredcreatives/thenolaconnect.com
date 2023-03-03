package model

import (
	"time"

	"gorm.io/gorm"
)

type PDFGenerationWorker struct {
	PDFGenerationWorkerId int            `json:"pdf_generation_worker_id" gorm:"column:pdf_generation_worker_id;autoIncrement;default:1"`
	FormId                string         `json:"form_id" gorm:"column:form_id"`
	OrderId               int            `json:"order_id" gorm:"column:order_id"`
	ProcessId             int            `json:"process_id" gorm:"column:process_id"`
	StartAt               time.Time      `json:"start_at" gorm:"column:start_at"`
	CompletedAt           time.Time      `json:"completed_at" gorm:"column:completed_at"`
	CreatedAt             time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (w *PDFGenerationWorker) TableName() string {
	return "order_communication.pdf_generation_worker"
}

func (w *PDFGenerationWorker) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&w)
}

func (w *PDFGenerationWorker) Delete(tx *gorm.DB) *gorm.DB {
	if w.PDFGenerationWorkerId > 0 {
		return tx.Delete(&w)
	}
	return tx
}

func (w *PDFGenerationWorker) FindOne(tx *gorm.DB) *gorm.DB {
	return tx.Model(&w).Find(&w)
}
