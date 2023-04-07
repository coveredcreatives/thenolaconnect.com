package communicator

import (
	"fmt"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
	"gorm.io/gorm"
)

type NewOrderCommunicator struct {
	OrderId int
	Db      *gorm.DB
}

func (c *NewOrderCommunicator) Order() (*model.Order, error) {
	var o model.Order
	if err := c.Db.First(&o, c.OrderId).Error; err != nil {
		return &o, err
	} else {
		return &o, nil
	}
}

func (c *NewOrderCommunicator) IsFulfilled(body string) (is_fulfilled bool, err error) {
	defer alog.Trace("NewOrderCommunicator.IsFulfilled").Stop(&err)
	o, err := c.Order()
	if err != nil {
		return
	}

	err = c.Db.Model(&model.Order{}).
		Select("count(*) > 0").
		Where("order_id = ? and first_conversation_message_sid != ?", o.OrderId, "").
		Find(&is_fulfilled).Error

	alog.WithField("is_fulfilled", is_fulfilled).WithField("order_id", o.OrderId).Info("is new order status fulfilled")
	return
}

func (c *NewOrderCommunicator) Store(body string) (is_valid bool, err error) {
	o, err := c.Order()
	if err != nil {
		return
	}
	is_valid = o.MediaSid != ""
	return
}

func (c *NewOrderCommunicator) Respond(body string) (params []conversations_openapi.CreateServiceConversationMessageParams, err error) {
	defer alog.Trace("NewOrderCommunicator.Respond").Stop(&err)
	o, err := c.Order()
	if err != nil {
		return
	}
	var message, media conversations_openapi.CreateServiceConversationMessageParams
	message.SetBody(fmt.Sprintf("Order #%05d is ready for manager review.\nReply (%s) to accept this order.\nReply (%s) to reject this order.", o.OrderId, AcceptOrderCode, RejectOrderCode))
	media.SetMediaSid(o.MediaSid)
	params = append(params, message, media)
	return
}
