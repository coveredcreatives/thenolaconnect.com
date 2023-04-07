package communicator

import (
	conversations_openapi "github.com/twilio/twilio-go/rest/conversations/v1"

	"github.com/coveredcreatives/thenolaconnect.com/model"
)

const AcceptOrderCode string = "1"
const RejectOrderCode string = "0"

type Communicator interface {
	Order() (*model.Order, error)
	IsFulfilled(body string) (bool, error)
	Store(body string) (bool, error)
	Respond(body string) ([]conversations_openapi.CreateServiceConversationMessageParams, error)
}
