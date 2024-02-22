package sms

type SmsInterface interface {

	Init()
	//返回短信网关中文名称
	GetName() string

	//返回接口中文介绍
	GetMemo() string

	//返回需要配置的项
	GetConfigInfo() []SmsConfigInfo

	//发送短信
	Send() SmsResult

	//获取短信余额
	Surplus() string

	//设置发信内容
	SetContent(content string)

	//获取发信内容
	GetContent() string

	//设置短信报告id
	SetSmsId(id string)

	//获取短信报告id
	GetSmsId() string

	//设置收件人
	SetAddress(address string)

	//获取收件人
	GetAddress() string

	//获取单条短信收费字数
	GetContentCharge() int64

	//返回发送失败的提示信息
	GetMessage()
}

type SmsConfigInfo struct {
	ConfigName  string
	ConfigValue string
	Name        string
	Type        string
	Style       string
}

type SmsResult struct {
	Status bool
	Info   string
	Used   int64
}
