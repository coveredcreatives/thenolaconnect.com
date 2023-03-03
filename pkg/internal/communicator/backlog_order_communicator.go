package communicator

import (
	"fmt"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
	"gorm.io/gorm"
)

type BacklogOrderCommunicator struct {
	OrderId int
	backlog []model.Order
	Db      *gorm.DB
}

func (c *BacklogOrderCommunicator) Order() (*model.Order, error) {
	var o model.Order
	if err := c.Db.First(&o, c.OrderId).Error; err != nil {
		return &o, err
	} else {
		return &o, nil
	}
}

func (c *BacklogOrderCommunicator) IsFulfilled(body string) (is_fulfilled bool, err error) {
	defer alog.Trace("BacklogOrderCommunicator.IsFulfilled").Stop(&err)
	var target []model.Order
	tx := c.Db.Model(&model.Order{}).
		Where("conversation_sid = ? AND is_viewed_by_manager = ?", "", false).
		Find(&target)
	err = tx.Error
	if err != nil {
		return
	}
	c.backlog = target
	is_fulfilled = tx.RowsAffected == 0
	alog.WithField("is_fulfilled", is_fulfilled).Info("is order validation fulfilled")
	return
}

func (c *BacklogOrderCommunicator) Store(body string) (is_valid bool, err error) {
	is_valid = true // incoming body is not used nor expected for backlog orders
	return
}

func (c *BacklogOrderCommunicator) Respond(body string) (params []conversations_openapi.CreateServiceConversationMessageParams, err error) {
	defer alog.Trace("BacklogOrderCommunicator.Respond").Stop(&err)
	var message conversations_openapi.CreateServiceConversationMessageParams
	message.SetBody(fmt.Sprintf("There are %d more orders to review.", len(c.backlog)))
	params = append(params, message)
	return
}
