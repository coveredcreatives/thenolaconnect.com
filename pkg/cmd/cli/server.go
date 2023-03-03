package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
	"gitlab.com/the-new-orleans-connection/qr-code/handlers"
	internal_tools "gitlab.com/the-new-orleans-connection/qr-code/internal"
	"google.golang.org/api/option"
)

func Server(ctx *cli.Context) error {
	googleauthconfig, err := devtools.GoogleApplicationConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing require google application configuration")
	}

	google_client_option := option.WithCredentialsFile(googleauthconfig.EnvGoogleApplicationCredentials)

	dbconfig, err := devtools.DatabaseConnectionConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing required database configuration")
		return err
	}
	gormdb, err := devtools.DatabaseConnection(ctx.Context, dbconfig)
	if err != nil {
		alog.WithError(err).Error("unable to open gorm instance")
		return err
	}

	storage_client, err := storage.NewClient(ctx.Context, google_client_option)
	if err != nil {
		alog.WithError(err).Error("unable to open google cloud storage client")
		return err
	}

	twilioconfig, err := devtools.TwilioConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing required twilio configuration")
		return err
	}

	twilio_client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioconfig.EnvAccountSid,
		Password: twilioconfig.EnvAccountAuthToken,
	})

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

	order_communication_config, err := devtools.OrderCommunicationConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing required database configuration")
		return err
	}

	alog.WithField("port", order_communication_config.EnvHTTPPort).Info("server listening")

	if err := http.ListenAndServe(fmt.Sprint(":", order_communication_config.EnvHTTPPort), router); err != nil {
		alog.WithError(err).Error(fmt.Sprint("unable to run server on port :", order_communication_config.EnvHTTPPort))
		return err
	}

	return nil
}
