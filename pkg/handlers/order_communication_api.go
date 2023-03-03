package handlers

import (
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"
)

func ListOrders(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		data, err := listOrders(w, r, db)
		if err != nil {
			_, _ = w.Write([]byte("bad request"))
		}
		_, _ = w.Write(data)
	})
}

func GenerateOrder(db *gorm.DB, storagec *storage.Client, twilioc *twilio.RestClient) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		data, err := generateOrder(w, r, db, storagec, twilioc)
		if err != nil {
			_, _ = w.Write([]byte("bad request"))
		}
		_, _ = w.Write(data)
	})
}

func DeliverOrderToKitchen(db *gorm.DB, twilioc *twilio.RestClient, orderchan chan<- int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		data, err := deliverOrderToKitchen(w, r, db, twilioc, orderchan)
		if err != nil {
			_, _ = w.Write([]byte("bad request"))
		}
		_, _ = w.Write(data)
	})
}

func SMS(db *gorm.DB, twilioc *twilio.RestClient, orderchan chan<- int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := sms(w, r, db, twilioc, orderchan)
		if err != nil {
			_, _ = w.Write([]byte("bad request"))
		}
		_, _ = w.Write(data)
	})
}
