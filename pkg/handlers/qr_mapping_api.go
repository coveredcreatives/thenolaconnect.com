package handlers

import (
	"image/png"
	"net/http"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"
)

func LoadV1QRMappingAPI(v *viper.Viper, router *http.ServeMux, gormdb *gorm.DB, storage_client *storage.Client, twilio_client *twilio.RestClient, orderchan chan int) {
	router.HandleFunc("/v1/qr_mapping/generate", Generate(gormdb, storage_client, v))
	router.HandleFunc("/v1/qr_mapping/list", Retrieve(gormdb, storage_client))
	router.HandleFunc("/v1/qr_mapping/retrieve", Retrieve(gormdb, storage_client))
}

func Retrieve(db *gorm.DB, storagec *storage.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		if r.URL.Query().Get("qr_encoded_data") != "" {
			redirect, err := retrieveOne(w, r, db, storagec)
			if err != nil {
				_, _ = w.Write([]byte("bad request"))
			} else {
				http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			}
		} else {
			data, err := retrieve(w, r, db, storagec)
			if err != nil {
				_, _ = w.Write([]byte("bad request"))
			}
			_, _ = w.Write(data)
		}
	})
}

// Generate QR handler will receive a file and a designated label.
// The file will be stored in a google cloud file storage bucket,
// an existing QR's metadata will be queried in the database based
// on its name matching the provided label. Matching QR is pulled from
// the bucket. If no match is found, a QR code will be created and uploaded
// the bucket.
func Generate(db *gorm.DB, storagec *storage.Client, v *viper.Viper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := alog.WithField("path", "pkg.handlers.generate.(GenerateQRHandler)")
		cors(w, r)
		entry.Info("BEGIN")
		qrCode := generate(w, r, db, storagec, v)
		_ = png.Encode(w, qrCode)
		entry.Info("END")
	})
}
