package devtools

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type DatabaseConnectionConfig struct {
	EnvDBIamUser              string `json:"DB_IAM_USER" env:"DB_IAM_USER"`
	EnvDBUsername             string `json:"DB_USERNAME" env:"DB_USERNAME,required"`
	EnvDBPassword             string `json:"DB_PASSWORD" env:"DB_PASSWORD"`
	EnvDBName                 string `json:"DB_NAME" env:"DB_NAME"`
	EnvDBPort                 int    `json:"DB_PORT" env:"DB_PORT,required"`
	EnvDBHostname             string `json:"DB_HOSTNAME" env:"DB_HOSTNAME,required"`
	EnvInstanceConnectionName string `json:"INSTANCE_CONNECTION_NAME" env:"INSTANCE_CONNECTION_NAME"`
	EnvInstanceUnixSocket     string `json:"UNIX_SOCKET_PATH" env:"UNIX_SOCKET_PATH"`
	EnvUsePrivate             string `json:"PRIVATE_IP" env:"PRIVATE_IP"`
}

type OrderCommunicationConfig struct {
	EnvName               string `json:"ENV" env:"ENV"`
	EnvHTTPPort           int    `json:"HTTP_PORT" env:"HTTP_PORT"`
	EnvPrinterIPv4Address string `json:"PRINTER_IPV4_ADDRESS" env:"PRINTER_IPV4_ADDRESS"`
}

type QRCodeConfig struct {
	EnvName                    string `json:"ENV" env:"ENV"`
	EnvHTTPPort                int    `json:"HTTP_PORT" env:"HTTP_PORT"`
	EnvRetrieveHTTPSTriggerUrl string `json:"RETRIEVE_HTTPS_TRIGGER_URL" env:"RETRIEVE_HTTPS_TRIGGER_URL"`
}

type GoogleApplicationConfig struct {
	EnvGoogleAPIKeyOrders                   string `json:"GOOGLE_API_KEY_ORDERS" env:"GOOGLE_API_KEY_ORDERS"`
	EnvGoogleFormIdOrders                   string `json:"GOOGLE_FORM_ID_ORDERS" env:"GOOGLE_FORM_ID_ORDERS"`
	EnvGoogleApplicationServiceAccountEmail string `json:"GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL" env:"GOOGLE_APPLICATION_SERVICE_ACCOUNT_EMAIL"`
	EnvGoogleApplicationCredentials         string `json:"GOOGLE_APPLICATION_CREDENTIALS" env:"GOOGLE_APPLICATION_CREDENTIALS"`
	EnvGoogleStorageBucketName              string `json:"GOOGLE_STORAGE_BUCKET_NAME" env:"GOOGLE_STORAGE_BUCKET_NAME"`
}

type TwilioConfig struct {
	EnvName                   string `json:"ENV" env:"ENV"`
	EnvAccountSid             string `json:"TWILIO_ACCOUNT_SID" env:"TWILIO_ACCOUNT_SID"`
	EnvConversationServiceSid string `json:"TWILIO_CONVERSATION_SERVICE_SID" env:"TWILIO_CONVERSATION_SERVICE_SID"`
	EnvMessagingServiceSid    string `json:"TWILIO_MESSAGING_SERVICE_SID" env:"TWILIO_MESSAGING_SERVICE_SID"`
	EnvAccountAuthToken       string `json:"TWILIO_ACCOUNT_AUTH_TOKEN" env:"TWILIO_ACCOUNT_AUTH_TOKEN"`
}

type PDFGenerationWorkerServerConfig struct {
	EnvBaseURL string `json:"PGW_BASE_URL" env:"PGW_BASE_URL,required"`
}

func DatabaseConnectionConfigFromEnv() (DatabaseConnectionConfig, error) {
	ctx := context.Background()

	var c DatabaseConnectionConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c, nil
}

func GoogleApplicationConfigFromEnv() (GoogleApplicationConfig, error) {
	ctx := context.Background()

	var c GoogleApplicationConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c, nil
}

func OrderCommunicationConfigFromEnv() (OrderCommunicationConfig, error) {
	ctx := context.Background()

	var c OrderCommunicationConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c, nil
}

func QRCodeConfigFromEnv() (QRCodeConfig, error) {
	ctx := context.Background()

	var c QRCodeConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c, nil
}

func TwilioConfigFromEnv() (TwilioConfig, error) {
	ctx := context.Background()

	var c TwilioConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c, nil
}

func PDFGenerationWorkerServerConfigFromEnv() (PDFGenerationWorkerServerConfig, error) {
	ctx := context.Background()

	var c PDFGenerationWorkerServerConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c, nil
}
