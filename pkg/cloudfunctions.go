package cloudfunctions

import (
	"context"
	"os"

	"cloud.google.com/go/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/handlers"
)

func init() {
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

	functions.HTTP("GenerateQRHandler", handlers.Generate(gormdb, storage_client, v))
	functions.HTTP("RetrieveQRHandler", handlers.Retrieve(gormdb, storage_client))
}
