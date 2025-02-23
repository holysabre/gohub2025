package sms

import (
	"gohub/pkg/config"
	"sync"
)

type Message struct {
	Template string
	Data     map[string]string
	Content  string
}

type SMS struct {
	Driver Driver
}

var once sync.Once

var intervalSMS *SMS

func NewSMS() *SMS {
	once.Do(func() {
		intervalSMS = &SMS{
			Driver: &Tencent{},
		}
	})
	return intervalSMS
}

func (sms *SMS) Send(phone string, message Message) bool {
	return sms.Driver.Send(phone, message, config.GetStringMapString("sms.tencent"))
}
