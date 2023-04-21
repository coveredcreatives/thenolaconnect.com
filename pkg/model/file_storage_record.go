package model

import (
	"time"

	"gorm.io/gorm"
)

type FileStorageRecord struct {
	FileStorageRecordId int            `json:"file_storage_record_id" gorm:"column:file_storage_record_id;primaryKey;autoincrement;default:1"`
	FileStorageUrl      string         `json:"file_storage_url" gorm:"column:file_storage_url"`
	CreatedAt           time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (f *FileStorageRecord) TableName() string {
	return "qr_mapping.file_storage_record"
}

func (f *FileStorageRecord) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&f)
}

func (f *FileStorageRecord) Delete(tx *gorm.DB) *gorm.DB {
	if f.FileStorageRecordId > 0 {
		return tx.Delete(&f)
	}
	return tx
}

func (f *FileStorageRecord) FindOne(tx *gorm.DB) *gorm.DB {
	return tx.Model(&f).Find(&f)
}
