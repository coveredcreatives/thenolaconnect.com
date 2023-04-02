package internal

import (
	"time"

	communication_tools "github.com/coveredcreatives/thenolaconnect.com/pkg/internal/communicator"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"
)

func ValidateOrder(v *viper.Viper, db *gorm.DB, twilioc *twilio.RestClient, order_id int, body string, orderchan chan<- int) (err error) {
	cc := communication_tools.CommunicationChain{}

	err = cc.
		Append(&communication_tools.ValidateOrderCommunicator{OrderId: order_id, Db: db, OrderChan: orderchan}).
		Run(v, db, twilioc, body)
	if err != nil {
		return
	}

	time.Sleep(time.Duration(1) * time.Second)
	_, _, err = CreateOrderConversation(v, db, twilioc)
	if err != gorm.ErrRecordNotFound && err != nil {
		return
	}
	return
}
