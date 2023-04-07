package internal

import (
	"os"
	"time"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	pgw_tools "github.com/coveredcreatives/thenolaconnect.com/internal/pgw"
	storage_tools "github.com/coveredcreatives/thenolaconnect.com/internal/storage"
	twilio_tools "github.com/coveredcreatives/thenolaconnect.com/internal/twilio"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func StreamPDFStoStorage(v *viper.Viper, db *gorm.DB, storagec *storage.Client, pgwchan <-chan model.PDFGenerationWorker, donechan chan<- bool) (err error) {
	defer alog.Trace("streamPDFStoStorage").Stop(&err)
	now := time.Now()

	google_env_config, err := devtools.GoogleApplicationLoadConfig(v)
	if err != nil {
		return
	}

	twilio_env_config, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		return
	}

	for pgw := range pgwchan {
		var b []byte
		b, err = pgw_tools.BuildOrderHTML(db, pgw)
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
