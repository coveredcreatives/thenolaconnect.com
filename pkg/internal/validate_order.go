package internal

import (
	"time"

	"github.com/twilio/twilio-go"
	communication_tools "gitlab.com/the-new-orleans-connection/qr-code/internal/communicator"
	"gorm.io/gorm"
)

func ValidateOrder(db *gorm.DB, twilioc *twilio.RestClient, order_id int, body string, orderchan chan<- int) (err error) {
	cc := communication_tools.CommunicationChain{}

	err = cc.
		Append(&communication_tools.ValidateOrderCommunicator{OrderId: order_id, Db: db, OrderChan: orderchan}).
		Run(db, twilioc, body)
	if err != nil {
		return
	}

	time.Sleep(time.Duration(1) * time.Second)
	_, _, err = CreateOrderConversation(db, twilioc)
	if err != gorm.ErrRecordNotFound && err != nil {
		return
	}
	return
}
