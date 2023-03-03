package internal

import (
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
)

func checkIsOrdersBacklogged(db *gorm.DB, twilioc *twilio.RestClient, phone_number string) (order model.Order, is_blocked bool, err error) {
	defer alog.WithField("phone_number", phone_number).Trace("checkIsOrdersBacklogged").Stop(&err)
	tx := db.
		Model(&model.Order{}).
		Where("conversation_sid != ? AND incoming_phone_number = ? AND is_viewed_by_manager = ?", "", phone_number, false).
		Find(&order)
	err = tx.Error
	if err != nil {
		return
	}
	is_blocked = tx.RowsAffected > 0

	alog.WithField("is_blocked", is_blocked).Info("find active orders using phone number from database")

	return
}
