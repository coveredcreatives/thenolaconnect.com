package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/handlers"
	internal_tools "github.com/coveredcreatives/thenolaconnect.com/pkg/internal"
	twilio "github.com/twilio/twilio-go"
)

func main() {

	dbconfig, err := devtools.DatabaseConnectionConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing required database configuration")
	}
	gormdb, err := devtools.DatabaseConnection(context.Background(), dbconfig)
	if err != nil {
		alog.WithError(err).Error("unable to open gorm instance")
	}

	storage_client, err := storage.NewClient(context.Background())
	if err != nil {
		alog.WithError(err).Error("unable to open google cloud storage client")
	}

	twilio_client := twilio.NewRestClient()

	orderchan := make(chan int)

	go internal_tools.ChannelOrdersToPrinter(gormdb, orderchan)

	router := http.NewServeMux()

	router.HandleFunc("/_ah/warmup", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("ok"))
	})

	router.HandleFunc("/qr_mapping/generate", handlers.Generate(gormdb, storage_client))
	router.HandleFunc("/qr_mapping/list", handlers.Retrieve(gormdb, storage_client))
	router.HandleFunc("/qr_mapping/retrieve", handlers.Retrieve(gormdb, storage_client))

	router.HandleFunc("/order_communication/generate", handlers.GenerateOrder(gormdb, storage_client, twilio_client))
	router.HandleFunc("/order_communication/list", handlers.ListOrders(gormdb))
	router.HandleFunc("/order_communication/deliver_to_kitchen", handlers.DeliverOrderToKitchen(gormdb, twilio_client, orderchan))
	router.HandleFunc("/order_communication/sms", handlers.SMS(gormdb, twilio_client, orderchan))

	qr_code_config, err := devtools.QRCodeConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing required database configuration")
	}

	// Support legacy lambda functions until rolled over
	functions.HTTP("GenerateQRHandler", handlers.Generate(gormdb, storage_client))
	functions.HTTP("RetrieveQRHandler", handlers.Retrieve(gormdb, storage_client))

	alog.WithField("port", qr_code_config.EnvHTTPPort).Info("server listening")

	if err := http.ListenAndServe(fmt.Sprint(":", qr_code_config.EnvHTTPPort), router); err != nil {
		alog.WithError(err).Error(fmt.Sprint("unable to run server on port :", qr_code_config.EnvHTTPPort))
	}
}
