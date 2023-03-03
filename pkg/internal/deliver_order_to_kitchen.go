package internal

import (
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"
)

func DeliverOrderToKitchen(db *gorm.DB, twilioc *twilio.RestClient, order_id int, orderchan chan<- int) {
	orderchan <- order_id
	alog.WithField("order_id", order_id).Info("deliver order to order channel")
}
