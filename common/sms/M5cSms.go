package sms

import (
	"fmt"
	"github.com/astaxie/beego"
	"html"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"shangchenggo/tools"
	"strings"
)

// +----------------------------------------------------------------------
// | 美联软通短信接口
// | api.52ao.com
// | M5cSms::send('手机号','内容');
// +----------------------------------------------------------------------
func (m *M5cSms) Init() {

	//读取数据库中的设置
	//$this->m5c_name = CONF('sms_M5cSms_name') ? CONF('sms_M5cSms_name') : 'MLRT短信接口';
	//$this->m5c_account = CONF('sms_M5cSms_account') ? CONF('sms_M5cSms_account') : '';
	//$this->m5c_passwd = CONF('sms_M5cSms_passwd') ? CONF('sms_M5cSms_passwd') : '';
	//$this->m5c_apikey = CONF('sms_M5cSms_apikey') ? CONF('sms_M5cSms_apikey') : '';
	//$this->sign_name = CONF('sms_M5cSms_sign') ? CONF('sms_M5cSms_sign')
	//测试帐号密码 key
	m.GatewayUrl = "http://m.5c.com.cn/api/send/index.php"
	m.M5cName = "MLRT短信接口"
	m.M5cAccount = "000008"
	m.M5cPasswd = "Duanxin013"
	m.M5cApikey = "22ddf4e3af5a29b73e22b7240f1d280b"
	m.SignName = "科技"
}

type M5cSms struct {
	GatewayUrl    string //发送网关地址
	M5cName       string
	M5cAccount    string
	M5cPasswd     string
	M5cApikey     string
	SignName      string
	Smsid         string
	Content       string
	Addressee     string
	ContentCharge int64
}

//返回接口中文名称 美联软通
func (m *M5cSms) GetName() string {
	return "MLRT"
}

//返回接口中文介绍
func (m *M5cSms) GetMemo() string {
	return "MLRT短信网关";
}

//返回需要配置的项
func (m *M5cSms) GetConfigInfo() []SmsConfigInfo {
	smsConfigInfos := []SmsConfigInfo{}
	smsc1 := SmsConfigInfo{
		ConfigName:  "sms_M5cSms_name",
		ConfigValue: "内置短信接口二",
		Name:        "内置短信接口",
		Type:        "text",
		Style:       "width:100px",
	}
	smsConfigInfos = append(smsConfigInfos, smsc1)

	smsc2 := SmsConfigInfo{
		ConfigName:  "sms_M5cSms_account",
		ConfigValue: "",
		Name:        "用户名",
		Type:        "text",
		Style:       "width:300px",
	}
	smsConfigInfos = append(smsConfigInfos, smsc2)

	smsc3 := SmsConfigInfo{
		ConfigName:  "sms_M5cSms_passwd",
		ConfigValue: "",
		Name:        "密码",
		Type:        "text",
		Style:       "width:380px",
	}
	smsConfigInfos = append(smsConfigInfos, smsc3)

	smsc4 := SmsConfigInfo{
		ConfigName:  "sms_M5cSms_apikey",
		ConfigValue: "",
		Name:        "API密钥",
		Type:        "text",
		Style:       "width:380px",
	}
	smsConfigInfos = append(smsConfigInfos, smsc4)

	smsc5 := SmsConfigInfo{
		ConfigName:  "sms_M5cSms_sign",
		ConfigValue: "",
		Name:        "签名",
		Type:        "text",
		Style:       "width:380px",
	}
	smsConfigInfos = append(smsConfigInfos, smsc5)
	return smsConfigInfos
}

/**
* 短信发送处理,根据网关类型
* @return array('status'=> true|false ,'info' => 返回信息)
*/
func (m *M5cSms) Send() SmsResult {
	encode := "utf8"
	user := m.M5cAccount
	content := fmt.Sprintf("【%s】%s", m.SignName, m.Content)
	apiKey := m.M5cApikey
	if "" != user && "" != m.M5cPasswd {
		gatewayurl := m.GatewayUrl
		mobiles := m.Addressee
		content = html.UnescapeString(content)
		pwd := strings.ToLower(tools.Md5(m.M5cPasswd))
		//url := fmt.Sprintf("%s?username=%s&password_md5=%s&apikey=%s&mobile=%s&content=%s&encode=%s",
		//	gatewayurl, user, pwd, apiKey, mobiles, content, encode)
		url, _ := url2.Parse(gatewayurl)
		values := url.Query()
		values.Set("username", user)
		values.Set("password_md5", pwd)
		values.Set("apikey", apiKey)
		values.Set("mobile", mobiles)
		values.Set("content", content)
		values.Set("encode", encode)
		url.RawQuery = values.Encode()
		beego.Info(url)
		by, err := getResult(url.String())
		if nil != err {
			beego.Error("getResult  err = ", err)
		}
		result := string(by)
		if strings.Contains(result, ":") {
			ress := strings.Split(result, ":")
			result = ress[0]
		}
		if "success" != result {
			ress := strings.Split(result, ":")
			errinfo := resolve_error(ress[1])
			return SmsResult{Status: false, Info: errinfo}
		} else {
			return SmsResult{Status: true, Info: "发送成功"}
		}
	}
	return SmsResult{Status: true, Info: "发送成功"}
}

/**
* 短信余额查看
* public $ddk_account  		= '';
* public $ddk_passwd  		= '';
* public $ddk_comcode 	= '';
*/
func (m *M5cSms) Surplus() string {
	user := m.M5cAccount
	if "" == user {
		return "-1"
	}
	pwd := tools.Md5(m.M5cPasswd)
	apiKey := m.M5cApikey
	url := fmt.Sprintf("http://m.5c.com.cn/api/query/index.php?username=%s&password_md5=%s&apikey=%s",
		user, pwd, apiKey)
	by, err := getResult(url)
	if nil != err {
		beego.Info("getResult err = ", err)
	}
	result := string(by)
	if !strings.Contains(result, "/") {
		ress := strings.Split(result, ":")
		return ress[1]
	} else {
		ress := strings.Split(result, "/")
		return ress[0]
	}
}

//设置发信内容
func (m *M5cSms) SetContent(content string) {
	m.Content = content
}

//获取发信内容
func (m *M5cSms) GetContent() string {
	return m.Content
}

//设置短信报告id
func (m *M5cSms) SetSmsId(id string) {
	m.Smsid = id
}

//获取短信报告id
func (m *M5cSms) GetSmsId() string {
	return m.Smsid
}

//设置收件人
func (m *M5cSms) SetAddress(address string) {
	m.Addressee = address
}

//获取收件人
func (m *M5cSms) GetAddress() string {
	return m.Addressee
}

//获取单条短信收费字数
func (m *M5cSms) GetContentCharge() int64 {
	return 200
}

//返回错误信息
func (m *M5cSms) GetMessage() {
}

func resolve_error(error_code string) string {
	msg := "未知"
	switch error_code {
	case "Missing username":
		msg = "用户名为空"
		break
	case "Missing password":
		msg = "密码为空"
		break
	case "Missing apikey":
		msg = "APIKEY为空"
		break
	case "Missing recipient":
		msg = "收件人手机号码为空"
		break
	case "Missing message content":
		msg = "短信内容为空或编码不正确"
		break
	case "Account is blocked":
		msg = "帐号被禁用"
		break
	case "Unrecognized encoding":
		msg = "编码未能识别"
		break

	case "APIKEY or password error":
		msg = "APIKEY或密码错误"
		break
	case "Unauthorized IP address":
		msg = "未授权 IP 地址"
		break

	case "Account balance is insufficient":
		msg = "余额不足"
		break

	case "Throughput Rate Exceeded":
		msg = "发送频率受限"
		break
	case "Invalid md5 password length":
		msg = "MD5密码长度非32位"
		break
	}
	return msg
}

func getResult(url string) ([]byte, error) {
	res, err := http.Get(url)
	if nil != err {
		beego.Error("http.NewRequest url = ", url)
	}
	if 200 != res.StatusCode {
		beego.Error("请求失败")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if nil != err {
		beego.Error("ioutil.ReadAll error")
	}
	return body, err
}
