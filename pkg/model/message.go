package model

type Message struct {
	Body                string `json:"body" gorm:"column:body"`
	NumSegments         string `json:"num_segments" gorm:"column:num_segments"`
	Direction           string `json:"direction" gorm:"column:direction"`
	From                string `json:"from" gorm:"column:from_pn"`
	To                  string `json:"to" gorm:"column:to_pn"`
	DateUpdated         string `json:"date_updated" gorm:"column:date_updated"`
	Price               string `json:"price" gorm:"column:price"`
	ErrorMessage        string `json:"error_message" gorm:"column:error_message"`
	Uri                 string `json:"uri" gorm:"column:uri"`
	AccountSid          string `json:"account_sid" gorm:"column:account_sid"`
	NumMedia            string `json:"num_media" gorm:"column:num_media"`
	Status              string `json:"status" gorm:"column:status"`
	MessagingServiceSid string `json:"messaging_service_sid" gorm:"column:messaging_service_sid"`
	Sid                 string `json:"sid" gorm:"column:sid"`
	DateSent            string `json:"date_sent" gorm:"column:date_sent"`
	DateCreated         string `json:"date_created" gorm:"column:date_created"`
	ErrorCode           int    `json:"error_code" gorm:"column:error_code"`
	PriceUnit           string `json:"price_unit" gorm:"column:price_unit"`
	ApiVersion          string `json:"api_version" gorm:"column:api_version"`
	SubresourceUris     string `json:"subresource_uris" gorm:"column:subresource_uris"`
}

func (m *Message) TableName() string {
	return "twilio.message"
}
