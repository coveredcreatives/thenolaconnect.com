package communicator

import (
	"fmt"

	alog "github.com/apex/log"
	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
)

type ValidateOrderCommunicator struct {
	OrderId   int
	Db        *gorm.DB
	OrderChan chan<- int
}

func (c *ValidateOrderCommunicator) Order() (*model.Order, error) {
	var o model.Order
	if err := c.Db.First(&o, c.OrderId).Error; err != nil {
		return &o, err
	} else {
		return &o, nil
	}
}

func (c *ValidateOrderCommunicator) IsFulfilled(body string) (is_fulfilled bool, err error) {
	defer alog.WithField("body", body).Trace("ValidateOrderCommunicator.IsFulfilled").Stop(&err)
	o, err := c.Order()
	if err != nil {
		return
	}
	is_fulfilled = o.IsViewedByManager
	alog.WithField("is_fulfilled", is_fulfilled).Info("is order validation fulfilled")
	return
}

func (c *ValidateOrderCommunicator) Store(body string) (is_valid bool, err error) {
	defer alog.WithField("body", body).Trace("ValidateOrderCommunicator.Store").Stop(&err)
	o, err := c.Order()
	if err != nil {
		return
	}
	if body == AcceptOrderCode {
		is_valid = true
		err = c.Db.Model(&model.Order{}).Where(&model.Order{OrderId: o.OrderId}).Updates(&model.Order{IsViewedByManager: true, IsAcceptedByManager: true}).Error
		c.OrderChan <- o.OrderId
		alog.WithField("order_id", o.OrderId).Info("deliver order to order channel")
	} else if body == RejectOrderCode {
		is_valid = true
		err = c.Db.Model(&model.Order{}).Where(&model.Order{OrderId: o.OrderId}).Updates(&model.Order{IsViewedByManager: true, IsAcceptedByManager: false}).Error
	} else {
		is_valid = false
	}
	alog.WithField("is_valid", is_valid).Info("is body valid to store")
	return
}

func (c *ValidateOrderCommunicator) Respond(body string) (params []conversations_openapi.CreateServiceConversationMessageParams, err error) {
	defer alog.Trace("ValidateOrderCommunicator.Respond").Stop(&err)
	o, err := c.Order()
	if err != nil {
		return
	}
	var message conversations_openapi.CreateServiceConversationMessageParams
	if o.IsViewedByManager && o.IsAcceptedByManager {
		message.SetBody(fmt.Sprintf("Order #%05d has been accepted, it is now queued up for kitchen delivery.", o.OrderId))
		params = append(params, message)
	}
	if o.IsViewedByManager && !o.IsAcceptedByManager {
		message.SetBody(fmt.Sprintf("Order #%05d has been rejected and must be resubmitted.", o.OrderId))
		params = append(params, message)
	}
	return
}
