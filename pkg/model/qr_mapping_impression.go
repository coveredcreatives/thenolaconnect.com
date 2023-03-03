package model

import (
	"time"

	"gorm.io/gorm"
)

type QRMappingImpression struct {
	QRMappingImpressionId int            `json:"qr_mapping_impression_id" gorm:"column:qr_mapping_impression_id;primaryKey;autoincrement;default:1"`
	QREncodedData         string         `json:"qr_encoded_data" gorm:"column:qr_encoded_data"`
	Host                  string         `json:"host" gorm:"column:host"`
	Path                  string         `json:"path" gorm:"column:path"`
	IPAddress             string         `json:"ip_address" gorm:"column:ip_address"`
	CreatedAt             time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (q *QRMappingImpression) TableName() string {
	return "qr_mapping.qr_mapping_impression"
}

func (q *QRMappingImpression) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&q)
}

func (q *QRMappingImpression) Delete(tx *gorm.DB) *gorm.DB {
	if q.QRMappingImpressionId > 0 {
		return tx.Delete(&q)
	}
	return tx
}

func (q *QRMappingImpression) FindOne(tx *gorm.DB) *gorm.DB {
	return tx.Model(&q).Find(&q)
}
