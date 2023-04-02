package internal

import (
	"fmt"

	alog "github.com/apex/log"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"

	"github.com/coveredcreatives/thenolaconnect.com/pkg/devtools"
	communication_tools "github.com/coveredcreatives/thenolaconnect.com/pkg/internal/communicator"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	"gorm.io/gorm"
)

func CreateOrderMessage(v *viper.Viper, db *gorm.DB, twilioc *twilio.RestClient) (err error) {
	defer alog.Trace("CreateOrderMessage").Stop(&err)

	twilio_env_config, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		return
	}

	var order model.Order
	var is_backlogged bool
	var incoming_pn_local model.IncomingPhoneNumberLocal

	tx := db.Model(&model.IncomingPhoneNumberLocal{}).First(&incoming_pn_local, model.IncomingPhoneNumberLocal{FriendlyName: fmt.Sprint("order_communication_", twilio_env_config.EnvName)})
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).Info("find phone number used for this application")

	cc := communication_tools.CommunicationChain{}

	order, is_backlogged, err = checkIsOrdersBacklogged(db, twilioc, incoming_pn_local.PhoneNumber)
	if err != nil {
		return
	}

	if !is_backlogged {
		order, err = findOrderWithNoConversationMessage(db)
	}

	if err != nil {
		return
	}

	err = cc.
		Append(&communication_tools.NewOrderCommunicator{Db: db, OrderId: order.OrderId}).
		Append(&communication_tools.BacklogOrderCommunicator{Db: db, OrderId: order.OrderId}).
		Run(v, db, twilioc, "")

	return
}

func findOrderWithNoConversationMessage(db *gorm.DB) (order model.Order, err error) {
	defer alog.Trace("findOrderWithNoConversationMessage").Stop(&err)

	tx := db.Model(&model.Order{}).
		Where("conversation_sid != ? AND conversation_sid NOT IN (?)", "", db.Model(&model.ConversationMessage{}).Select("conversation_sid")).
		Limit(1).
		Find(&order)
	err = tx.Error
	alog.WithField("rows_affected", tx.RowsAffected).Info("find orders in conversation with no conversation_messages from database")

	return
}

func IdentifyOrderFromPhoneNumberLocal(db *gorm.DB, service_phone string) (order model.Order, err error) {
	defer alog.WithField("service_phone", service_phone).Trace("identifyOrderFromPhoneNumberLocal").Stop(&err)

	tx := db.Model(&model.Order{}).Where("incoming_phone_number = ? AND is_viewed_by_manager = ?", service_phone, false).Find(&order)
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("num_rows", tx.RowsAffected).WithField("order_id", order.OrderId).Info("find orders with incoming phone number from database")
	return
}
