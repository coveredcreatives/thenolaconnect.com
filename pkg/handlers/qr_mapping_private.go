package handlers

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
	storageTools "gitlab.com/the-new-orleans-connection/qr-code/internal/storage"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func retrieve(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client) ([]byte, error) {
	alog.Info("begin retrieve")
	qr_mappings := []model.QRMapping{}
	db.Model(&qr_mappings).Find(&qr_mappings)
	qr_mapping_pkeys := make([]string, len(qr_mappings))
	alog.WithField("mappings", len(qr_mappings)).Info("fetched qr mappings")
	for i, qr_mapping := range qr_mappings {
		qr_mapping_pkeys[i] = qr_mapping.QREncodedData
	}
	qr_mapping_metas_by_id := model.CalculateUniqueImpressions(db, qr_mapping_pkeys)
	response := []map[string]interface{}{}
	for _, qr_mapping := range qr_mappings {
		alog.WithField("qr_encoded_data", qr_mapping.QREncodedData).Info("calculated unique impressions")
		response = append(response, map[string]interface{}{
			"qr_encoded_data":        qr_mapping.QREncodedData,
			"qr_file_storage_url":    qr_mapping.QRFileStorageUrl,
			"name":                   qr_mapping.Name,
			"file_storage_record_id": qr_mapping.FileStorageRecordId,
			"unique_impressions":     qr_mapping_metas_by_id[qr_mapping.QREncodedData].UniqueImpressions,
		})
	}
	return json.Marshal(response)
}

func retrieveOne(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client) (string, error) {
	alog.Info("begin retrieve one")
	qrEncodedData := r.URL.Query().Get("qr_encoded_data")
	alog.WithField("data", qrEncodedData).Info("new retrieve request")
	qr_mapping := model.QRMapping{}
	filestoragerecord := model.FileStorageRecord{}
	_ = db.First(&qr_mapping, model.QRMapping{QREncodedData: qrEncodedData})
	alog.WithField("name", qr_mapping.Name).WithField("name", qr_mapping.Name).WithField("url", qr_mapping.QRFileStorageUrl).Info("fetched qr mapping")
	_ = db.First(&filestoragerecord, model.FileStorageRecord{FileStorageRecordId: qr_mapping.FileStorageRecordId})
	alog.WithField("filestorage_id", filestoragerecord.FileStorageRecordId).WithField("url", filestoragerecord.FileStoragedUrl).Info("fetched file storage record")
	return filestoragerecord.FileStoragedUrl, nil
}

func generate(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client) barcode.Barcode {
	entry := alog.WithField("path", "pkg.handlers.generate.(generate)")
	entry.Info("BEGIN")
	googleapplicationconfig, err := devtools.GoogleApplicationConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("failed to find Google Application config")
		_, _ = w.Write([]byte("bad request"))
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		entry.WithError(err).Error("failed to reference file from request")
		_, _ = w.Write([]byte("bad request"))
	}
	entry.WithField("filename", header.Filename).WithField("filesize", header.Size).Info("received file from request")
	objAttrs, err := storageTools.UploadFileToStorage(storagec, file, header, googleapplicationconfig.EnvGoogleStorageBucketName)
	if err != nil {
		entry.WithError(err).Error("failed to upload file to storage")
		_, _ = w.Write([]byte("bad request"))
	}
	entry.WithField("name", objAttrs.Name).WithField("link", storageTools.FmtPublicUrl(objAttrs)).Info("uploaded file to storage")

	filestoragerecord := &model.FileStorageRecord{FileStoragedUrl: storageTools.FmtPublicUrl(objAttrs)}
	result := db.Clauses(clause.Returning{}).Create(filestoragerecord)
	if result.Error != nil {
		entry.WithError(err).Error("failed to store file record in db")
		_, _ = w.Write([]byte("bad request"))
	}

	entry.WithField("fileStorageRecordId", filestoragerecord.FileStorageRecordId).Info("successfully stored record in db")

	label := strings.ReplaceAll(strings.ToLower(r.FormValue("label")), " ", "-")
	qrMapping := &model.QRMapping{}
	qrencodeddata := b64.URLEncoding.EncodeToString([]byte(label))

	qrcodeconfig, err := devtools.QRCodeConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("failed to find QR mapping")
		_, _ = w.Write([]byte("bad request"))
	}

	retrieveURL := fmt.Sprintf("https://%s/retrieve?qr_encoded_data=%s", qrcodeconfig.EnvRetrieveHTTPSTriggerUrl, qrencodeddata)
	alog.WithField("url", retrieveURL).Info("stored retrieveURL for code")
	qrCode, _ := qr.Encode(retrieveURL, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	qrObjAttrs, err := storageTools.UploadImageToStorage(storagec, qrCode, fmt.Sprintf("QR-%s", qrMapping.QREncodedData), googleapplicationconfig.EnvGoogleStorageBucketName)
	if err != nil {
		alog.WithError(err).Error("failed to find upload qrCode to file storage")
		_, _ = w.Write([]byte("bad request"))
	}
	qrmapping := model.QRMapping{
		Name:                label,
		QREncodedData:       qrencodeddata,
		QRFileStorageUrl:    storageTools.FmtPublicUrl(qrObjAttrs),
		FileStorageRecordId: filestoragerecord.FileStorageRecordId,
	}
	result = db.
		Create(&qrmapping)
	if result.Error != nil {
		alog.WithError(err).Error("failed to find QR mapping")
		_, _ = w.Write([]byte("bad request"))
	}
	alog.WithField("qr_mapping", qrMapping).Info("created new qr mapping")

	return qrCode
}
