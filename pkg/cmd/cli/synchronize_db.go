package main

import (
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	internal_tools "github.com/coveredcreatives/thenolaconnect.com/internal"
	"github.com/spf13/viper"

	"github.com/twilio/twilio-go"
	cli "github.com/urfave/cli/v2"

	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

func SynchronizeDB(ctx *cli.Context, v *viper.Viper) error {
	googleauthconfig, err := devtools.GoogleApplicationLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing require google application configuration")
	}

	google_client_option := option.WithCredentialsFile(googleauthconfig.EnvGoogleApplicationCredentials)

	forms_service, err := forms.NewService(ctx.Context, google_client_option)
	if err != nil {
		alog.WithError(err).Error("unable to initialize new service")
	}

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

	twilio_env_config, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing required twilio configuration")
		return err
	}

	twilioc := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilio_env_config.EnvAccountSid,
		Password: twilio_env_config.EnvAccountAuthToken,
	})

	err = internal_tools.SynchronizeDB(v, forms_service, gormdb, twilioc)
	if err != nil {
		return err
	}

	return nil
}
