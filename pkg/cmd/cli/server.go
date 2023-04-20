package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/handlers"
	internal_tools "github.com/coveredcreatives/thenolaconnect.com/internal"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

func Server(ctx *cli.Context, v *viper.Viper) error {
	googleauthconfig, err := devtools.GoogleApplicationLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing require google application configuration")
	}

	google_client_option := option.WithCredentialsFile(googleauthconfig.EnvGoogleApplicationCredentials)

	dbconfig, err := devtools.DatabaseConnectionLoadConfig(v)
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

	twilioconfig, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing required twilio configuration")
		return err
	}

	twilio_client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioconfig.EnvAccountSid,
		Password: twilioconfig.EnvAccountAuthToken,
	})

	forms_service, err := forms.NewService(context.Background())
	if err != nil {
		alog.WithError(err).Error("unable to initialize new service")
	}

	orderchan := make(chan int)

	go internal_tools.ChannelOrdersToPrinter(v, gormdb, orderchan)

	router := http.NewServeMux()

	router.HandleFunc("/_ah/warmup", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("ok"))
	})

	handlers.LoadV1OrderCommunicationAPI(v, router, gormdb, storage_client, twilio_client, forms_service, orderchan)
	handlers.LoadV1QRMappingAPI(v, router, gormdb, storage_client, twilio_client, orderchan)

	application_config, err := devtools.ApplicationLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing required database configuration")
		return err
	}

	alog.WithField("port", application_config.EnvHTTPPort).Info("server listening")

	if err := http.ListenAndServe(fmt.Sprint(":", application_config.EnvHTTPPort), router); err != nil {
		alog.WithError(err).Error(fmt.Sprint("unable to run server on port :", application_config.EnvHTTPPort))
		return err
	}

	return nil
}
