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
	router.HandleFunc("/v1/qr_mapping/generate", Generate(v, gormdb, storage_client))
	router.HandleFunc("/v1/qr_mapping/list", Retrieve(v, gormdb, storage_client))
	router.HandleFunc("/v1/qr_mapping/retrieve", Retrieve(v, gormdb, storage_client))
	router.HandleFunc("/v1/qr_mapping/hide", Hide(v, gormdb, storage_client))
	router.HandleFunc("/v1/qr_mapping/show", Show(v, gormdb, storage_client))
}

func Retrieve(v *viper.Viper, db *gorm.DB, storagec *storage.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := alog.WithField("path", "pkg.handlers.qr_mapping_api.(Retrieve)")
		cors(v, w, r)
		entry.Info("BEGIN")
		if r.URL.Query().Get("qr_encoded_data") != "" {
			redirect, err := retrieveOne(w, r, db, storagec)
			if err != nil {
				entry.WithError(err).Error("END")
				_, _ = w.Write([]byte("bad request"))
				return
			} else {
				http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			}
		} else {
			data, err := retrieve(w, r, db, storagec)
			if err != nil {
				entry.WithError(err).Error("END")
				_, _ = w.Write([]byte("bad request"))
				return
			}
			_, _ = w.Write(data)
			entry.Info("END")
			return
		}
	})
}

func Generate(v *viper.Viper, db *gorm.DB, storagec *storage.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := alog.WithField("path", "pkg.handlers.qr_mapping_api.(Generate)")
		cors(v, w, r)
		entry.Info("BEGIN")
		qrCode, err := generate(w, r, db, storagec, v)
		if err != nil {
			entry.WithError(err).Error("END")
			return
		}
		err = png.Encode(w, qrCode)
		if err != nil {
			entry.WithError(err).Error("END")
			return
		}
		entry.Info("END")
	})
}

func Hide(v *viper.Viper, db *gorm.DB, storagec *storage.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := alog.WithField("path", "pkg.handlers.qr_mapping_api.(Hide)")
		cors(v, w, r)
		entry.Info("BEGIN")
		ok, err := hide(w, r, db, storagec, v)
		if err != nil {
			entry.WithError(err).Error("END")
			return
		}
		if ok {
			w.WriteHeader(http.StatusNoContent)
			entry.Info("END")
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			entry.Info("END")
			return
		}
	})
}

func Show(v *viper.Viper, db *gorm.DB, storagec *storage.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := alog.WithField("path", "pkg.handlers.qr_mapping_api.(Show)")
		cors(v, w, r)
		entry.Info("BEGIN")
		ok, err := show(w, r, db, storagec, v)
		if err != nil {
			entry.WithError(err).Error("END")
			return
		}
		if ok {
			w.WriteHeader(http.StatusNoContent)
			entry.Info("END")
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			entry.Info("END")
			return
		}
	})
}
