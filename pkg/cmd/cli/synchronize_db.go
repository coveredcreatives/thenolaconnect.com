package main

import (
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
	form_tools "gitlab.com/the-new-orleans-connection/qr-code/internal/forms"
	twilio_tools "gitlab.com/the-new-orleans-connection/qr-code/internal/twilio"

	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

func SynchronizeDB(ctx *cli.Context) error {
	googleauthconfig, err := devtools.GoogleApplicationConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing require google application configuration")
	}

	google_client_option := option.WithCredentialsFile(googleauthconfig.EnvGoogleApplicationCredentials)

	forms_service, err := forms.NewService(ctx.Context, google_client_option)
	if err != nil {
		alog.WithError(err).Error("unable to initialize new service")
	}

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

	twilio_env_config, err := devtools.TwilioConfigFromEnv()
	if err != nil {
		alog.WithError(err).Error("missing required twilio configuration")
		return err
	}

	twilioc := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilio_env_config.EnvAccountSid,
		Password: twilio_env_config.EnvAccountAuthToken,
	})

	form_tools.ListFormSynchronizeDB(gormdb, forms_service, googleauthconfig.EnvGoogleFormIdOrders)

	if form_tools.ListFormResponsesSynchronizeDB(gormdb, forms_service, googleauthconfig.EnvGoogleFormIdOrders, "") != nil {
		return err
	}

	if twilio_tools.ListPhoneNumbersSynchronizeDB(gormdb, twilioc, twilio_env_config.EnvAccountSid) != nil {
		return err
	}

	if twilio_tools.FetchAccountSynchronizeToDB(gormdb, twilioc, twilio_env_config.EnvAccountSid) != nil {
		return err
	}

	if twilio_tools.FetchServiceSynchronizeDB(gormdb, twilioc, twilio_env_config.EnvConversationServiceSid, twilio_env_config.EnvMessagingServiceSid) != nil {
		return err
	}

	return nil
}
