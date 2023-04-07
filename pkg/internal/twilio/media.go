package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	"gorm.io/gorm"
)

func StoreTwilioMedia(db *gorm.DB, order model.Order, twilio_env_config devtools.TwilioConfig) (media model.Media, err error) {
	defer alog.Trace("StoreTwilioMedia").Stop(&err)
	client := http.Client{}
	filename := order.MediaFilename()
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	defer os.Remove(filename)

	request, err := http.NewRequest(http.MethodPost, fmt.Sprint("https://mcs.us1.twilio.com/v1/Services/", twilio_env_config.EnvConversationServiceSid, "/Media"), f)
	if err != nil {
		return
	}

	request.SetBasicAuth(twilio_env_config.EnvAccountSid, twilio_env_config.EnvAccountAuthToken)
	request.Header.Add("Content-Type", "application/pdf")
	values := request.URL.Query()
	values.Add("Filename", fmt.Sprintf("order_%05d.pdf", order.OrderId))
	request.URL.RawQuery = values.Encode()

	response, err := client.Do(request)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &media)
	alog.WithField("order.OrderId", order.OrderId).WithField("media_url", media.Url).WithField("media_sid", media.Sid).Info("media uploaded")
	return
}
