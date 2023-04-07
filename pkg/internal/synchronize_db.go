package internal

import (
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	form_tools "github.com/coveredcreatives/thenolaconnect.com/internal/forms"
	twilio_tools "github.com/coveredcreatives/thenolaconnect.com/internal/twilio"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"

	"google.golang.org/api/forms/v1"
)

func SynchronizeDB(v *viper.Viper, forms_service *forms.Service, gormdb *gorm.DB, twilioc *twilio.RestClient) error {
	googleauthconfig, err := devtools.GoogleApplicationLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing require google application configuration")
	}
	twilio_env_config, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		alog.WithError(err).Error("missing required twilio configuration")
		return err
	}
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
