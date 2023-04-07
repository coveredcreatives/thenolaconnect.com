package internal

import (
	"encoding/json"
	"fmt"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"
	"gorm.io/gorm"
)

func CreateOrderConversation(v *viper.Viper, db *gorm.DB, twilioc *twilio.RestClient) (order model.Order, is_blocked bool, err error) {
	defer alog.Trace("createOrderConversation").Stop(&err)

	twilio_env_config, err := devtools.TwilioLoadConfig(v)
	if err != nil {
		return
	}

	var incoming_pn_local model.IncomingPhoneNumberLocal
	tx := db.Model(&model.IncomingPhoneNumberLocal{}).First(&incoming_pn_local, model.IncomingPhoneNumberLocal{FriendlyName: fmt.Sprint("order_communication_", twilio_env_config.EnvName)})
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).Info("find phone number used for this application")

	order, is_blocked, err = checkIsOrdersBacklogged(db, twilioc, incoming_pn_local.PhoneNumber)
	if err != nil {
		return
	}
	if is_blocked {
		return
	}

	var conversations_v1_conversations_by_form_id map[string]conversations_openapi.ConversationsV1ServiceConversation
	conversations_v1_conversations_by_form_id, err = findAndGroupTwilioServiceConversationList(twilioc, twilio_env_config.EnvConversationServiceSid)
	managers := []model.Staff{}
	tx = db.
		Model(&model.Staff{}).
		Where(&model.Staff{
			IsManager: true,
		}).
		Find(&managers)
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).Info("find staff from database")

	tx = db.Model(&model.Order{}).Where("conversation_sid = ? AND is_viewed_by_manager = ?", "", false).First(&order)
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("order.OrderId", order.OrderId).Info("recieve order from order channel")

	conversation, ok := conversations_v1_conversations_by_form_id[order.FormId]
	if !ok {
		conversation, err = createServiceConversation(db, twilioc, order, twilio_env_config.EnvConversationServiceSid)
		if err != nil {
			return
		}
	}

	tx = db.
		Model(&model.Order{}).
		Where(&model.Order{OrderId: order.OrderId}).
		Updates(&model.Order{
			ConversationSid:     *conversation.Sid,
			IncomingPhoneNumber: incoming_pn_local.PhoneNumber,
		})
	err = tx.Error
	if err != nil {
		return
	}
	alog.Info("update order with conversation sid, phone number in database")

	var c_participant_by_message_binding_address map[string]conversations_openapi.ConversationsV1ServiceConversationParticipant
	c_participant_by_message_binding_address, err = findAndGroupTwilioServiceConversationParticpant(twilioc, twilio_env_config.EnvConversationServiceSid, *conversation.Sid)

	for _, m := range managers {
		var rows_affected int64
		var is_associated_to_participant bool
		manager := m
		tx := db.Model(&model.ConversationParticipant{}).Limit(1).Where(&model.ConversationParticipant{
			MessagingBindingAddress: manager.PhoneNumber,
			ConversationSid:         *conversation.Sid,
		}).Count(&rows_affected)
		err = tx.Error
		if err != nil {
			return
		}
		alog.WithField("conversation_sid", *conversation.Sid).
			WithField("message_binding_address", manager.PhoneNumber).
			WithField("rows_affected", rows_affected).
			Info("find conversation participant from database")
		is_associated_to_participant = rows_affected == 1
		if is_associated_to_participant {
			return
		}

		is_associated_to_participant, err = associateExistingServiceConversationParticipant(db, c_participant_by_message_binding_address, manager.PhoneNumber)
		if err != nil {
			return
		}
		if is_associated_to_participant {
			return
		}

		_, err = createServiceConversationParticipant(db, twilioc, twilio_env_config.EnvConversationServiceSid, *conversation.Sid, manager.PhoneNumber, incoming_pn_local.PhoneNumber)
	}
	return
}

func findAndGroupTwilioServiceConversationList(twilioc *twilio.RestClient, conversation_service_sid string) (conversations_v1_conversations_by_response_id map[string]conversations_openapi.ConversationsV1ServiceConversation, err error) {
	defer alog.WithField("conversation_service_sid", conversation_service_sid).Trace("findAndGroupTwilioServiceConversationList").Stop(&err)
	conversations_v1_conversations, err := twilioc.ConversationsV1.ListServiceConversation(conversation_service_sid, &conversations_openapi.ListServiceConversationParams{})
	if err != nil {
		return
	}
	conversations_v1_conversations_by_response_id = make(map[string]conversations_openapi.ConversationsV1ServiceConversation, len(conversations_v1_conversations))
	for _, c := range conversations_v1_conversations {
		conversations_v1_conversations_by_response_id[*c.UniqueName] = c
	}
	return
}

func findAndGroupTwilioServiceConversationParticpant(twilioc *twilio.RestClient, conversation_service_sid, conversation_sid string) (c_participant_by_message_binding_address map[string]conversations_openapi.ConversationsV1ServiceConversationParticipant, err error) {
	defer alog.WithField("conversation_service_sid", conversation_service_sid).WithField("conversation_sid", conversation_sid).Trace("findAndGroupTwilioServiceConversationParticipant").Stop(&err)
	var conversations_v1_c_participants []conversations_openapi.ConversationsV1ServiceConversationParticipant
	conversations_v1_c_participants, err = twilioc.
		ConversationsV1.
		ListServiceConversationParticipant(conversation_service_sid, conversation_sid, &conversations_openapi.ListServiceConversationParticipantParams{})
	if err != nil {
		return
	}
	alog.WithField("num_rows", len(conversations_v1_c_participants)).Info("list service conversation participants from twilio")
	c_participant_by_message_binding_address = make(map[string]conversations_openapi.ConversationsV1ServiceConversationParticipant, len(conversations_v1_c_participants))
	for _, c_participant := range conversations_v1_c_participants {
		messaging_binding := map[string]string{"address": "", "proxy_address": ""}
		binding, _ := json.Marshal(c_participant.MessagingBinding)
		_ = json.Unmarshal(binding, &messaging_binding)
		c_participant_by_message_binding_address[messaging_binding["address"]] = c_participant
	}
	return
}

func createServiceConversation(db *gorm.DB, twilioc *twilio.RestClient, order model.Order, conversation_service_sid string) (conversation conversations_openapi.ConversationsV1ServiceConversation, err error) {
	defer alog.WithField("order.OrderId", order.OrderId).WithField("conversation_service_sid", conversation_service_sid).Trace("createServiceConversation").Stop(&err)
	var conversations_v1_conversation *conversations_openapi.ConversationsV1ServiceConversation
	twilio_conversation_params := conversations_openapi.CreateServiceConversationParams{}
	twilio_conversation_params.SetFriendlyName(fmt.Sprint("Verifying Order #", order.OrderId))
	twilio_conversation_params.SetUniqueName(order.FormId)
	twilio_conversation_params.SetTimersInactive("P30D")
	twilio_conversation_params.SetTimersClosed("PT0S")
	conversations_v1_conversation, err = twilioc.ConversationsV1.CreateServiceConversation(conversation_service_sid, &twilio_conversation_params)
	if err != nil {
		return
	}
	alog.WithField("conversations_v1_conversation.Sid", conversations_v1_conversation.Sid).Info("create service conversation in twilio")

	c := model.ConversationToLocalSchema(*conversations_v1_conversation)
	tx := db.Model(&model.Conversation{}).Create(&c)
	err = tx.Error
	if err != nil {
		return
	}
	alog.Info("create conversation in database")
	conversation = *conversations_v1_conversation
	return
}

func associateExistingServiceConversationParticipant(db *gorm.DB, c_participant_by_message_binding_address map[string]conversations_openapi.ConversationsV1ServiceConversationParticipant, manager_pn string) (ok bool, err error) {
	defer alog.WithField("manger_phone", manager_pn).Trace("associateExistingServiceConversationparticipant").Stop(&err)

	conversations_v1_c_participant, ok := c_participant_by_message_binding_address[manager_pn]
	if !ok {
		return
	}
	c_participant := model.ConversationParticipantToSchema(conversations_v1_c_participant)

	err = db.
		Model(&model.ConversationParticipant{}).
		Create(&c_participant).Error
	if err != nil {
		return
	}

	alog.WithField("ok", ok).Info("create conversation participant in database")

	return
}

func createServiceConversationParticipant(db *gorm.DB, twilioc *twilio.RestClient, conversation_service_sid, conversation_sid, manager_pn, incoming_pn_local string) (c_participant model.ConversationParticipant, err error) {
	defer alog.WithField("conversation_service_sid", conversation_service_sid).
		WithField("conversation_sid", conversation_sid).
		WithField("manager_pn", manager_pn).
		WithField("incoming_pn_local", incoming_pn_local).
		Trace("createServiceConversationParticipants").Stop(&err)

	twilio_participant_params := conversations_openapi.CreateServiceConversationParticipantParams{}
	twilio_participant_params.SetMessagingBindingAddress(manager_pn)
	twilio_participant_params.SetMessagingBindingProxyAddress(incoming_pn_local)
	conversations_v1_c_participant, err := twilioc.
		ConversationsV1.
		CreateServiceConversationParticipant(conversation_service_sid, conversation_sid, &twilio_participant_params)
	if err != nil {
		return
	}
	alog.WithField("sid", conversations_v1_c_participant.Sid).WithField("message_binding", conversations_v1_c_participant.MessagingBinding).Info("create service conversation participant in twilio")

	c_participant = model.ConversationParticipantToSchema(*conversations_v1_c_participant)

	tx := db.
		Model(&model.ConversationParticipant{}).
		Create(&c_participant)

	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).Info("create participant in database")
	return
}
