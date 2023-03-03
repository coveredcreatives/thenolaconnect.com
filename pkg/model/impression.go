package model

import (
	"time"

	"gorm.io/gorm"
)

type Impression struct {
	ImpressionId int            `json:"impression_id" gorm:"column:impression_id;primaryKey;autoincrement;default:1"`
	Path         string         `json:"path" gorm:"column:path"`
	IPAddress    string         `json:"ip_address" gorm:"column:ip_address"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (i *Impression) TableName() string {
	return "qr_mapping.impression"
}

func (i *Impression) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&i)
}

func (i *Impression) Delete(tx *gorm.DB) *gorm.DB {
	if i.ImpressionId > 0 {
		return tx.Delete(&i)
	}
	return tx
}

func (i *Impression) FindOne(tx *gorm.DB) *gorm.DB {
	return tx.Model(&i).Find(&i)
}
