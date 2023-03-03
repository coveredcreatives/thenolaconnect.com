package model

import (
	api_openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type IncomingPhoneNumberLocal struct {
	AccountSid        string `json:"account_sid" gorm:"column:account_sid"`
	AddressSid        string `json:"address_sid" gorm:"column:address_sid"`
	FriendlyName      string `json:"friendly_name" gorm:"column:friendly_name"`
	IdentitySid       string `json:"identity_sid" gorm:"column:identity_sid"`
	PhoneNumber       string `json:"phone_number" gorm:"column:phone_number"`
	Sid               string `json:"sid" gorm:"column:sid"`
	SmsApplicationSid string `json:"sms_application_sid" gorm:"column:sms_application_sid"`
	Uri               string `json:"uri" gorm:"column:uri"`
	BundleSid         string `json:"bundle_sid" gorm:"column:bundle_sid"`
	Status            string `json:"status" gorm:"column:status"`
}

func (p *IncomingPhoneNumberLocal) TableName() string {
	return "twilio.incoming_phone_number_local"
}

func IncomingPhoneNumberLocalToLocalSchema(pn []api_openapi.ApiV2010IncomingPhoneNumberLocal) []IncomingPhoneNumberLocal {
	out := []IncomingPhoneNumberLocal{}
	for _, p := range pn {
		m := IncomingPhoneNumberLocal{}
		if p.AccountSid != nil {
			m.AccountSid = *p.AccountSid
		}
		if p.AddressSid != nil {
			m.AddressSid = *p.AddressSid
		}
		if p.FriendlyName != nil {
			m.FriendlyName = *p.FriendlyName
		}
		if p.IdentitySid != nil {
			m.IdentitySid = *p.IdentitySid
		}
		if p.PhoneNumber != nil {
			m.PhoneNumber = *p.PhoneNumber
		}
		if p.Sid != nil {
			m.Sid = *p.Sid
		}
		if p.SmsApplicationSid != nil {
			m.SmsApplicationSid = *p.SmsApplicationSid
		}
		if p.Uri != nil {
			m.Uri = *p.Uri
		}
		if p.BundleSid != nil {
			m.BundleSid = *p.BundleSid
		}
		if p.Status != nil {
			m.Status = *p.Status
		}
		out = append(out, m)
	}
	return out
}
