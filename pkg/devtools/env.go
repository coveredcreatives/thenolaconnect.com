package devtools

import (
	"os"

	"github.com/spf13/viper"
)

func Config() (v *viper.Viper, err error) {
	v = viper.New()

	v.SetEnvPrefix("nola")
	v.BindEnv("ENV")
	v.BindEnv("DB_USERNAME")
	v.BindEnv("DB_PASSWORD")
	v.BindEnv("DB_NAME")
	v.BindEnv("DB_PORT")
	v.BindEnv("DB_HOSTNAME")
	v.BindEnv("APP_CONFIG_PATH")
	v.BindEnv("HTTP_PORT")
	v.BindEnv("DNS_PRINTER_IPV4_ADDRESS")
	v.BindEnv("DNS_RETRIEVE_TRIGGER_URL")
	v.BindEnv("GOOGLE_API_KEY_ORDERS")
	v.BindEnv("GOOGLE_FORM_ID_ORDERS")
	v.BindEnv("GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL")
	v.BindEnv("GOOGLE_APPLICATION_CREDENTIALS")
	v.BindEnv("GOOGLE_STORAGE_BUCKET_NAME")
	v.BindEnv("TWILIO_ACCOUNT_SID")
	v.BindEnv("TWILIO_CONVERSATION_SERVICE_SID")
	v.BindEnv("TWILIO_MESSAGING_SERVICE_SID")
	v.BindEnv("TWILIO_ACCOUNT_AUTH_TOKEN")

	// Set Default Values where appropriate
	v.SetDefault("HTTP_PORT", 3001)
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	v.SetDefault("APP_CONFIG_FILENAME", v.GetString("ENV"))
	v.SetDefault("APP_CONFIG_PATH", dir)

	// Loads configuration file from ./appconfig.{yaml|json}
	v.SetConfigName(v.GetString("APP_CONFIG_FILENAME"))
	v.AddConfigPath(v.GetString("APP_CONFIG_PATH"))
	_ = v.ReadInConfig()
	return
}

type DatabaseConnectionConfig struct {
	EnvDBUsername string `json:"DB_USERNAME" mapstructure:"DB_USERNAME"`
	EnvDBPassword string `json:"DB_PASSWORD" mapstructure:"DB_PASSWORD"`
	EnvDBName     string `json:"DB_NAME" mapstructure:"DB_NAME"`
	EnvDBPort     int    `json:"DB_PORT" mapstructure:"DB_PORT"`
	EnvDBHostname string `json:"DB_HOSTNAME" mapstructure:"DB_HOSTNAME"`
}

type ApplicationConfig struct {
	EnvName                    string `json:"ENV" mapstructure:"ENV"`
	EnvHTTPPort                int    `json:"HTTP_PORT" mapstructure:"HTTP_PORT"`
	EnvPrinterIPv4Address      string `json:"PRINTER_IPV4_ADDRESS" mapstructure:"PRINTER_IPV4_ADDRESS"`
	EnvRetrieveHTTPSTriggerUrl string `json:"RETRIEVE_HTTPS_TRIGGER_URL" mapstructure:"RETRIEVE_HTTPS_TRIGGER_URL"`
}

type GoogleApplicationConfig struct {
	EnvGoogleFormIdOrders                   string `json:"GOOGLE_FORM_ID_ORDERS" mapstructure:"GOOGLE_FORM_ID_ORDERS"`
	EnvGoogleApplicationServiceAccountEmail string `json:"GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL" mapstructure:"GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL"`
	EnvGoogleApplicationCredentials         string `json:"GOOGLE_APPLICATION_CREDENTIALS" mapstructure:"GOOGLE_APPLICATION_CREDENTIALS"`
	EnvGoogleStorageBucketName              string `json:"GOOGLE_STORAGE_BUCKET_NAME" mapstructure:"GOOGLE_STORAGE_BUCKET_NAME"`
}

type TwilioConfig struct {
	EnvName                   string `json:"ENV" mapstructure:"ENV"`
	EnvAccountSid             string `json:"TWILIO_ACCOUNT_SID" mapstructure:"TWILIO_ACCOUNT_SID"`
	EnvConversationServiceSid string `json:"TWILIO_CONVERSATION_SERVICE_SID" mapstructure:"TWILIO_CONVERSATION_SERVICE_SID"`
	EnvMessagingServiceSid    string `json:"TWILIO_MESSAGING_SERVICE_SID" mapstructure:"TWILIO_MESSAGING_SERVICE_SID"`
	EnvAccountAuthToken       string `json:"TWILIO_ACCOUNT_AUTH_TOKEN" mapstructure:"TWILIO_ACCOUNT_AUTH_TOKEN"`
}

func DatabaseConnectionLoadConfig(v *viper.Viper) (DatabaseConnectionConfig, error) {

	var c DatabaseConnectionConfig

	v.Unmarshal(&c)

	return c, nil
}

func GoogleApplicationLoadConfig(v *viper.Viper) (GoogleApplicationConfig, error) {
	var c GoogleApplicationConfig

	v.Unmarshal(&c)

	return c, nil
}

func ApplicationLoadConfig(v *viper.Viper) (ApplicationConfig, error) {
	var c ApplicationConfig

	v.Unmarshal(&c)

	return c, nil
}

func TwilioLoadConfig(v *viper.Viper) (TwilioConfig, error) {
	var c TwilioConfig

	v.Unmarshal(&c)

	return c, nil
}
