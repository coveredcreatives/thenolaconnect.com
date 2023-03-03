package template

import (
	"gitlab.com/the-new-orleans-connection/qr-code/model"
)

type Page struct {
	Title              string                    `json:"title"`
	QRMappings         []model.QRMapping         `json:"qrMappings"`
	FileStorageRecords []model.FileStorageRecord `json:"fileStorageRecords"`
}
