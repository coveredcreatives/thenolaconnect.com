package twilio

import (
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
)

func FetchAccountSynchronizeToDB(db *gorm.DB, twilioc *twilio.RestClient, account_sid string) error {
	var err error
	tracer := alog.Trace("FetchAccountSynchronizeToDB")
	defer tracer.Stop(&err)

	api_v2010_account, err := twilioc.
		Api.
		FetchAccount(account_sid)
	if err != nil {
		tracer.Stop(&err)
		return err
	} else {
		tracer.Stop(nil)
	}

	acct := model.AccountToLocalSchema(*api_v2010_account)
	tracer = alog.Trace("create/update latest account information in database")
	if tx := db.
		Model(&model.Account{}).
		FirstOrCreate(&acct, &model.Account{
			Sid: account_sid,
		}); tx.Error != nil {
		tracer.Stop(&err)
		return err
	} else {
		tracer.Stop(nil)
	}

	return nil
}
