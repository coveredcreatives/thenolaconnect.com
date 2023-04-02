package handlers

import (
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	"google.golang.org/api/forms/v1"
	"gorm.io/gorm"
)

func LoadV1OrderCommunicationAPI(v *viper.Viper, router *http.ServeMux, gormdb *gorm.DB, storage_client *storage.Client, twilio_client *twilio.RestClient, forms_service *forms.Service, orderchan chan int) {
	router.HandleFunc("/v1/order_communication/generate", GenerateOrder(v, gormdb, storage_client, twilio_client))
	router.HandleFunc("/v1/order_communication/list", ListOrders(gormdb))
	router.HandleFunc("/v1/order_communication/deliver_to_kitchen", DeliverOrderToKitchen(gormdb, twilio_client, orderchan))
	router.HandleFunc("/v1/order_communication/sms", SMS(v, gormdb, twilio_client, orderchan))
	router.HandleFunc("/v1/order_communication/sync", SynchronizeDB(v, forms_service, gormdb, twilio_client))

}

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

func GenerateOrder(v *viper.Viper, db *gorm.DB, storagec *storage.Client, twilioc *twilio.RestClient) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		data, err := generateOrder(v, w, r, db, storagec, twilioc)
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

func SMS(v *viper.Viper, db *gorm.DB, twilioc *twilio.RestClient, orderchan chan<- int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := sms(v, w, r, db, twilioc, orderchan)
		if err != nil {
			_, _ = w.Write([]byte("bad request"))
		}
		_, _ = w.Write(data)
	})
}

func SynchronizeDB(v *viper.Viper, forms_service *forms.Service, db *gorm.DB, twilioc *twilio.RestClient) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := synchronizeDB(v, w, r, forms_service, db, twilioc)
		if err != nil {
			_, _ = w.Write([]byte("bad request"))
		}
		_, _ = w.Write(data)
	})
}
