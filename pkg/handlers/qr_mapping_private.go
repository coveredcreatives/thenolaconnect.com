package handlers

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	storageTools "github.com/coveredcreatives/thenolaconnect.com/internal/storage"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func retrieve(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client) (out []byte, err error) {
	entry := alog.WithField("path", "pkg.handlers.qr_mapping_private.(retrieve)")
	entry.Info("BEGIN")
	qr_mappings := []model.QRMapping{}
	db.Model(&qr_mappings).Find(&qr_mappings)
	qr_mapping_pkeys := make([]string, len(qr_mappings))
	entry.WithField("mappings", len(qr_mappings)).Info("fetched qr mappings")
	for i, qr_mapping := range qr_mappings {
		qr_mapping_pkeys[i] = qr_mapping.QREncodedData
	}
	qr_mapping_metas_by_id := model.CalculateUniqueImpressions(db, qr_mapping_pkeys)
	entry.WithField("qr_encoded_data", qr_mapping_pkeys).Info("calculated unique impressions")
	response := []map[string]interface{}{}
	for _, qr_mapping := range qr_mappings {
		response = append(response, map[string]interface{}{
			"qr_encoded_data":        qr_mapping.QREncodedData,
			"qr_file_storage_url":    qr_mapping.QRFileStorageUrl,
			"name":                   qr_mapping.Name,
			"file_storage_record_id": qr_mapping.FileStorageRecordId,
			"unique_impressions":     qr_mapping_metas_by_id[qr_mapping.QREncodedData].UniqueImpressions,
		})
	}
	out, err = json.Marshal(response)
	entry.Info("END")
	return
}

func retrieveOne(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client) (string, error) {
	entry := alog.WithField("path", "pkg.handlers.qr_mapping_private.(retrieveOne)")
	entry.Info("BEGIN")
	qrEncodedData := r.URL.Query().Get("qr_encoded_data")
	entry.WithField("data", qrEncodedData).Info("new retrieve request")
	qr_mapping := model.QRMapping{}
	filestoragerecord := model.FileStorageRecord{}
	_ = db.First(&qr_mapping, model.QRMapping{QREncodedData: qrEncodedData})
	entry.WithField("name", qr_mapping.Name).WithField("name", qr_mapping.Name).WithField("url", qr_mapping.QRFileStorageUrl).Info("fetched qr mapping")
	_ = db.First(&filestoragerecord, model.FileStorageRecord{FileStorageRecordId: qr_mapping.FileStorageRecordId})
	entry.WithField("filestorage_id", filestoragerecord.FileStorageRecordId).WithField("url", filestoragerecord.FileStorageUrl).Info("fetched file storage record")
	entry.Info("END")
	return filestoragerecord.FileStorageUrl, nil
}

func generate(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client, v *viper.Viper) (qrCode barcode.Barcode, err error) {
	entry := alog.WithField("path", "pkg.handlers.qr_mapping_private.(generate)")
	entry.Info("BEGIN")

	file, header, err := r.FormFile("file")
	if err != nil {
		entry.WithError(err).Error("failed to reference file from request")
		_, _ = w.Write([]byte("bad request"))
		return
	}
	entry.WithField("filename", header.Filename).WithField("filesize", header.Size).Info("received file from request")
	objAttrs, err := storageTools.UploadFileToStorage(storagec, file, header, v.GetString("GOOGLE_STORAGE_BUCKET_NAME"))
	if err != nil {
		entry.WithError(err).Error("failed to upload file to storage")
		_, _ = w.Write([]byte("bad request"))
		return
	}
	entry.WithField("name", objAttrs.Name).WithField("link", storageTools.FmtPublicUrl(objAttrs)).Info("uploaded file to storage")

	if db.Model(&model.FileStorageRecord{}).Clauses(clause.Returning{}).Create(&map[string]interface{}{
		"FileStorageRecordId": clause.Expr{SQL: "nextval('qr_mapping.file_storage_record_file_storage_record_id_seq'::regclass)"},
		"FileStorageUrl":      storageTools.FmtPublicUrl(objAttrs),
		"CreatedAt":           time.Now(),
		"UpdatedAt":           time.Now(),
	}).Error != nil {
		entry.WithError(err).Error("failed to store file record in db")
		_, _ = w.Write([]byte("bad request"))
		return
	}
	filestoragerecord := model.FileStorageRecord{}
	if db.Model(&model.FileStorageRecord{}).Where("file_storage_url = ?", storageTools.FmtPublicUrl(objAttrs)).Last(&filestoragerecord).Error != nil {
		entry.WithError(err).Error("failed to fetch file record id from db")
		_, _ = w.Write([]byte("bad request"))
		return
	}

	label := strings.ReplaceAll(strings.ToLower(r.FormValue("label")), " ", "-")
	qrMapping := &model.QRMapping{}
	qrencodeddata := b64.URLEncoding.EncodeToString([]byte(label))

	retrieveURL := fmt.Sprint("http://", r.Host, "/v1/qr_mapping/retrieve?qr_encoded_data=", qrencodeddata)
	alog.WithField("url", retrieveURL).Info("stored retrieveURL for code")
	qrCode, err = qr.Encode(retrieveURL, qr.L, qr.Auto)
	if err != nil {
		return
	}
	qrCode, _ = barcode.Scale(qrCode, 512, 512)
	if err != nil {
		return
	}

	qrObjAttrs, err := storageTools.UploadImageToStorage(storagec, qrCode, fmt.Sprintf("QR-%s", qrMapping.QREncodedData), v.GetString("GOOGLE_STORAGE_BUCKET_NAME"))
	if err != nil {
		alog.WithError(err).Error("failed to find upload qrCode to file storage")
		_, _ = w.Write([]byte("bad request"))
		return
	}
	var count int64
	tx := db.Model(&model.QRMapping{}).Where("qr_encoded_data = ?", qrencodeddata).Count(&count)
	if tx.Error != nil {
		alog.WithError(tx.Error).Error("failed to count from db")
		_, _ = w.Write([]byte("bad request"))
		return
	}
	if count > 0 {
		result := db.Model(&model.QRMapping{}).
			Where("qr_encoded_data = ?", qrencodeddata).
			Update("file_storage_record_id", filestoragerecord.FileStorageRecordId)
		if result.Error != nil {
			alog.WithError(err).Error("failed update mapping with new file storage record id")
			_, _ = w.Write([]byte("bad request"))
			return
		}
		alog.Info("updated qr mapping")
	} else {
		result := db.Create(&model.QRMapping{
			Name:                label,
			QREncodedData:       qrencodeddata,
			QRFileStorageUrl:    storageTools.FmtPublicUrl(qrObjAttrs),
			FileStorageRecordId: filestoragerecord.FileStorageRecordId,
		})
		if result.Error != nil {
			alog.WithError(err).Error("failed to create QR mapping")
			_, _ = w.Write([]byte("bad request"))
			return
		}
		alog.Info("created new qr mapping")
	}

	return
}

func hide(w http.ResponseWriter, r *http.Request, db *gorm.DB, storagec *storage.Client, v *viper.Viper) (ok bool, err error) {
	entry := alog.WithField("path", "pkg.handlers.qr_mapping_private.(hide)")
	entry.Info("BEGIN")

	if err = db.Model(&model.QRMapping{}).
		Where("qr_encoded_data = ?", r.URL.Query().Get("qr_encoded_data")).
		Update("deleted_at", time.Now()).Error; err != nil {
		entry.WithError(err).Error("failed to delete record")
		return
	}

	entry.Info("END")
	return
}

func show(w http.ResponseWriter, r *http.Request, db *gorm.DB, storagec *storage.Client, v *viper.Viper) (ok bool, err error) {
	entry := alog.WithField("path", "pkg.handlers.qr_mapping_private.(show)")
	entry.Info("BEGIN")

	if err = db.Model(&model.QRMapping{}).
		Where("qr_encoded_data = ?", r.URL.Query().Get("qr_encoded_data")).
		Update("deleted_at", nil).Error; err != nil {
		entry.WithError(err).Error("failed to show record")
		return
	}

	entry.Info("END")
	return
}
