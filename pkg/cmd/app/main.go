package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/forms/v1"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/handlers"
	internal_tools "github.com/coveredcreatives/thenolaconnect.com/pkg/internal"
	twilio "github.com/twilio/twilio-go"
)

func main() {
	v, err := devtools.Config()
	if err != nil {
		alog.WithError(err).Error("failed to load application configuration values")
		os.Exit(1)
	}

	dbconfig, err := devtools.DatabaseConnectionLoadConfig(v)
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

	forms_service, err := forms.NewService(context.Background())
	if err != nil {
		alog.WithError(err).Error("unable to initialize new service")
	}

	// Long running proccess to queue orders to be printed
	orderchan := make(chan int)

	go internal_tools.ChannelOrdersToPrinter(v, gormdb, orderchan)

	// Begin routing http service
	router := http.NewServeMux()

	router.HandleFunc("/_ah/warmup", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("ok"))
	})

	handlers.LoadV1OrderCommunicationAPI(v, router, gormdb, storage_client, twilio_client, forms_service, orderchan)
	handlers.LoadV1QRMappingAPI(v, router, gormdb, storage_client, twilio_client, orderchan)

	application_config, err := devtools.ApplicationLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing required database configuration")
	}

	alog.WithField("port", application_config.EnvHTTPPort).Info("server listening")

	if err := http.ListenAndServe(fmt.Sprint(":", application_config.EnvHTTPPort), router); err != nil {
		alog.WithError(err).Error(fmt.Sprint("unable to run server on port :", application_config.EnvHTTPPort))
	}
}
