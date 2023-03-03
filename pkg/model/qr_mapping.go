package model

import (
	"time"

	"gorm.io/gorm"
)

type QRMapping struct {
	FileStorageRecordId int            `json:"file_storage_record_id" gorm:"column:file_storage_record_id;default:1"`
	QRFileStorageUrl    string         `json:"qr_file_storage_url" gorm:"column:qr_file_storage_url"`
	Name                string         `json:"name" gorm:"column:name"`
	QREncodedData       string         `json:"qr_encoded_data" gorm:"column:qr_encoded_data;primaryKey;not null"`
	CreatedAt           time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (q *QRMapping) TableName() string {
	return "qr_mapping.qr_mapping"
}

func (q *QRMapping) Create(tx *gorm.DB) *gorm.DB {
	if q.QREncodedData != "" {
		return tx.Create(&q)
	}
	return tx
}

func (q *QRMapping) Delete(tx *gorm.DB) *gorm.DB {
	if q.QREncodedData != "" {
		return tx.Delete(&q)
	}
	return tx
}

func (q *QRMapping) FindOne(tx *gorm.DB) *gorm.DB {
	return tx.Model(&q).First(&q)
}
