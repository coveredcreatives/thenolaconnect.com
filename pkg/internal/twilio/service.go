package twilio

import (
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"
)

func FetchServiceSynchronizeDB(db *gorm.DB, twilioc *twilio.RestClient, conversation_service_sid string, messaging_service_sid string) error {
	var err error
	tracer := alog.Trace("FetchServiceSynchronizeDB")
	defer tracer.Stop(&err)

	twilio_api_tracer := alog.Trace("executing fetch conversations service")
	conversations_v1_service, err := twilioc.
		ConversationsV1.
		FetchService(conversation_service_sid)
	if err != nil {
		twilio_api_tracer.Stop(&err)
		return err
	} else {
		twilio_api_tracer.Stop(nil)
	}

	conversation_service := model.ConversationServiceToLocalSchema(*conversations_v1_service)
	db_tracer := alog.Trace("create/update twilio.conversation_service in db")
	if tx := db.
		Model(&model.ConversationService{}).
		FirstOrCreate(&conversation_service, &model.ConversationService{
			Sid: conversation_service_sid,
		}); tx.Error != nil {
		db_tracer.Stop(&err)
		return err
	} else {
		db_tracer.Stop(nil)
	}

	twilio_api_tracer = alog.Trace("executing fetch messaging service")
	messaging_v1_service, err := twilioc.MessagingV1.FetchService(messaging_service_sid)
	if err != nil {
		twilio_api_tracer.Stop(&err)
		return err
	} else {
		twilio_api_tracer.Stop(nil)
	}

	service := model.ServiceToLocalSchema(*messaging_v1_service)
	db_tracer = alog.Trace("create/update twilio.service in db")
	if tx := db.
		Model(&model.Service{}).
		FirstOrCreate(&service, &model.Service{
			Sid: messaging_service_sid,
		}); tx.Error != nil {
		db_tracer.Stop(&err)
		return err
	} else {
		db_tracer.Stop(nil)
	}

	return nil
}
