package communicator

import (
	"reflect"
	"time"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
	"gorm.io/gorm"
)

type CommunicationChain struct {
	c    Communicator
	next *CommunicationChain
}

func (cc *CommunicationChain) Append(c Communicator) *CommunicationChain {
	if cc.c == nil {
		cc.c = c
		alog.WithField("Communicator.Name", reflect.TypeOf(c).Elem().Name()).Info("Append")
	} else {
		if cc.next == nil {
			cc.next = &CommunicationChain{}
		}
		cc.next.Append(c)
	}
	return cc
}

func (c *CommunicationChain) Run(v *viper.Viper, db *gorm.DB, twilioc *twilio.RestClient, body string) (err error) {
	defer alog.WithField("Communicator.Name", reflect.TypeOf(c.c).Elem().Name()).Trace("Run").Stop(&err)

	twilio_env_config, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		return
	}

	order, err := c.c.Order()
	if err != nil {
		return
	}

	is_fulfilled, err := c.c.IsFulfilled(body)
	if err != nil {
		return
	}

	if is_fulfilled {
		err = traverse(c, v, db, twilioc, body)
		return
	}

	store_ok, err := storeMessage(c, body)
	if err != nil {
		return
	}

	twilio_message_params, err := calculateResponse(c, store_ok, body)
	if err != nil {
		return
	}

	c_message, err := deliverResponseToTwilio(twilioc, c, twilio_env_config.EnvConversationServiceSid, order.ConversationSid, twilio_message_params)
	if err != nil {
		return
	}

	err = db.Model(&model.ConversationMessage{}).Create(&c_message).Error
	if err != nil {
		return
	}

	err = db.Model(&model.Order{OrderId: order.OrderId}).Updates(&model.Order{FirstConversationMessageSid: c_message[0].Sid}).Error
	if err != nil {
		return
	}

	time.Sleep(time.Duration(1) * time.Second)
	err = traverse(c, v, db, twilioc, body)
	return
}

func storeMessage(c *CommunicationChain, body string) (ok bool, err error) {
	defer alog.WithField("body", body).Trace("storeMessage").Stop(&err)
	ok, err = c.c.Store(body)
	if err != nil {
		return
	}
	alog.WithField("ok", ok).Info("store incoming message")
	return
}

func calculateResponse(c *CommunicationChain, ok bool, body string) (twilio_message_params []conversations_openapi.CreateServiceConversationMessageParams, err error) {
	defer alog.Trace("responseMessage").Stop(&err)
	if ok {
		twilio_message_params, err = c.c.Respond(body)
		if err != nil {
			return
		}
	} else {
		var param conversations_openapi.CreateServiceConversationMessageParams
		param.SetBody("Sorry, we did not recognize your response.")
		twilio_message_params = append(twilio_message_params, param)
	}
	return
}

func deliverResponseToTwilio(twilioc *twilio.RestClient, c *CommunicationChain, conversation_service_sid string, conversation_sid string, twilio_message_params []conversations_openapi.CreateServiceConversationMessageParams) (c_message []model.ConversationMessage, err error) {
	defer alog.Trace("deliverResponseToTwilio").Stop(&err)
	var conversations_v1_c_message = new(conversations_openapi.ConversationsV1ServiceConversationMessage)
	for _, param := range twilio_message_params {
		conversations_v1_c_message, err = twilioc.ConversationsV1.CreateServiceConversationMessage(conversation_service_sid, conversation_sid, &param)
		if err != nil {
			return
		}
		c_message = append(c_message, model.ConversationMessageToSchema(*conversations_v1_c_message))
	}

	return
}

func traverse(c *CommunicationChain, v *viper.Viper, db *gorm.DB, twilioc *twilio.RestClient, body string) (err error) {
	if c.next != nil {
		err = c.next.Run(v, db, twilioc, body)
		return
	}
	return
}
