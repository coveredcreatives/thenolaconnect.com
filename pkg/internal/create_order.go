package internal

import (
	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	"github.com/twilio/twilio-go"

	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateOrder(db *gorm.DB, storagec *storage.Client, twilioc *twilio.RestClient, responsechan <-chan model.FormResponse, orderchan chan<- model.Order, errchan chan<- error) {
	defer alog.Trace("CreateOrder").Stop(nil)
	for r := range responsechan {
		response := r
		alog.WithField("response_id", response.ResponseId).Info("receive from channel")
		var order model.Order
		order.FormResponseId = response.ResponseId
		order.FormId = response.FormId
		tx := db.
			Clauses(clause.Returning{Columns: []clause.Column{{Table: "order_communication.order", Name: "order_id"}}}).
			Omit("order_id").
			Create(&order)
		err := tx.Error
		if err != nil {
			errchan <- err
			return
		}
		alog.
			WithField("rows_affected", tx.RowsAffected).
			WithField("order.OrderId", order.OrderId).
			WithField("order.FormId", order.FormId).
			WithField("order.FormResponseId", order.FormResponseId).
			Info("create order from response")
		orderchan <- order
		alog.WithField("order.OrderId", order.OrderId).Info("send order to order channel")
	}
	close(orderchan)
}
