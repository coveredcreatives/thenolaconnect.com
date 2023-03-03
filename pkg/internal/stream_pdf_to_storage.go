package internal

import (
	"os"
	"time"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
	pgw_tools "gitlab.com/the-new-orleans-connection/qr-code/internal/pgw"
	storage_tools "gitlab.com/the-new-orleans-connection/qr-code/internal/storage"
	twilio_tools "gitlab.com/the-new-orleans-connection/qr-code/internal/twilio"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
)

func StreamPDFStoStorage(db *gorm.DB, storagec *storage.Client, pgwchan <-chan model.PDFGenerationWorker, donechan chan<- bool) (err error) {
	defer alog.Trace("streamPDFStoStorage").Stop(&err)
	now := time.Now()

	pgw_env_config, err := devtools.PDFGenerationWorkerServerConfigFromEnv()
	if err != nil {
		return
	}

	google_env_config, err := devtools.GoogleApplicationConfigFromEnv()
	if err != nil {
		return
	}

	twilio_env_config, err := devtools.TwilioConfigFromEnv()
	if err != nil {
		return
	}

	for pgw := range pgwchan {
		var b []byte
		b, err = pgw_tools.TriggerPDFConversionApi(pgw_env_config, pgw)
		if err != nil {
			return
		}
		if len(b) == 0 {
			continue
		}

		order := model.Order{}
		err = db.First(&order, model.Order{OrderId: pgw.OrderId}).Error
		if err != nil {
			return
		}

		filename := order.MediaFilename()

		err = os.WriteFile(filename, b, 0644)
		if err != nil {
			return
		}

		storage_tools.StoreCloudStorageObject(db, storagec, pgw, order, b, google_env_config, now)

		if order.MediaSid == "" {
			media := model.Media{}
			media, err = twilio_tools.StoreTwilioMedia(db, order, twilio_env_config)
			if err != nil {
				return
			}

			err = db.Create(&media).Error
			if err != nil {
				return
			}

			err = db.Model(&model.Order{}).Where(&model.Order{OrderId: order.OrderId}).Updates(model.Order{MediaSid: media.Sid}).Error
			if err != nil {
				return
			}
		}

		os.Remove(filename)
	}
	donechan <- true
	return
}
