package template

import (
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
)

type Page struct {
	Title              string                    `json:"title"`
	QRMappings         []model.QRMapping         `json:"qrMappings"`
	FileStorageRecords []model.FileStorageRecord `json:"fileStorageRecords"`
}
