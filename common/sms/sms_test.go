package sms

import (
	"shangchenggo/conf"
	"testing"
)

func Test_Send(t *testing.T) {
	conf.Init()
	sms := new(Sms)
	sms.Prepare("M5cSms","hallow。。。。","")
	sms.Send()
}
