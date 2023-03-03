package model

import (
	"gorm.io/gorm"
)

type ICreate interface {
	TableName() string
	Create(*gorm.DB) *gorm.DB
}

type IDelete interface {
	TableName() string
	Delete(*gorm.DB) *gorm.DB
}

func Tables(tables []string) map[string]interface{} {
	index := map[string]interface{}{
		"file_storage_record":   &FileStorageRecord{},
		"impression":            &Impression{},
		"qr_mapping_impression": &QRMappingImpression{},
		"qr_mapping_meta":       &QRMappingMeta{},
		"qr_mapping":            &QRMapping{},
		"order":                 &Order{},
		"pdf_generation_worker": &PDFGenerationWorker{},
		"twilio.conversation":   &Conversation{},
	}
	if len(tables) > 0 {
		out := map[string]interface{}{}
		for _, table := range tables {
			if match, ok := index[table]; ok {
				out[table] = match
			}
		}
		return out
	} else {
		return index
	}
}
