package twilio

import (
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"
	accounts_openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
)

func ListPhoneNumbersSynchronizeDB(db *gorm.DB, twilioc *twilio.RestClient, account_sid string) error {
	var err error
	tracer := alog.Trace("ListPhoneNumbersSynchronizeDB")
	defer tracer.Stop(&err)

	twilio_api_tracer := alog.WithField("path_account_sid", account_sid).Trace("fetching incoming phone number local from twilio")
	api_v2_list_incoming_pn, err := twilioc.Api.ListIncomingPhoneNumberLocal(&accounts_openapi.ListIncomingPhoneNumberLocalParams{
		PathAccountSid: &account_sid,
	})
	if err != nil {
		twilio_api_tracer.Stop(&err)
		return err
	} else {
		twilio_api_tracer.Stop(nil)
	}

	for _, pn := range model.IncomingPhoneNumberLocalToLocalSchema(api_v2_list_incoming_pn) {
		db_create_tracer := alog.Trace("fill phone records from api into db")
		if tx := db.
			Omit("Capabilities").
			Model(&model.IncomingPhoneNumberLocal{}).
			FirstOrCreate(&pn, &model.IncomingPhoneNumberLocal{
				Sid:         pn.Sid,
				AccountSid:  pn.AccountSid,
				PhoneNumber: pn.PhoneNumber,
			}); tx.Error != nil {
			db_create_tracer.Stop(&err)
			return err
		} else {
			db_create_tracer.Stop(nil)
		}
	}

	return nil
}
