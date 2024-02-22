package sms

import (
	"github.com/astaxie/beego"
	"math"
	"shangchenggo/common/db"
	"shangchenggo/models"
	"shangchenggo/tools"
	"strconv"
	"strings"
	"time"
)

var (
	smsMap     map[string]interface{}
	smsapi     SmsInterface
	Smsid      string
	Content    string
	Address    string
	batSendNum int64
	errorMsgs  string
)

type Sms struct {
	smsname string
	content string
	address string
}

/**
	初始化
 */
func (s *Sms) Prepare(smsname string, content string, address string) {
	s.smsname = smsname
	smsapi = smsMap[smsname].(SmsInterface)
	smsapi.Init()
	batSendNum = 150
	newAddress := []string{}
	if strings.Contains(address, ",") {
		for _, v := range strings.Split(address, ",") {
			//过滤空手机号码
			if len(v) <= 0 {
				continue
			}
			newAddress = append(newAddress, v)
		}
	} else {
		newAddress = append(newAddress, address)
	}
	//再次转换成按 逗号 分隔的字符串格式
	Address = strings.Join(newAddress, ",")
	s.address = Address

	//产生短信编号
	Smsid = createSmsId()
	Content = content
	if len(smsname) > 0 && checkSmsName(smsname) {
		smsapi = smsMap[smsname].(SmsInterface)
		smsapi.SetContent(Content)
		smsapi.SetAddress(Address)
		smsapi.SetSmsId(Smsid)
		if len(Content) > 0 {
			createReport("", smsname)
		}
	} else {
		printMessage("请指定要使用的短信网关名称!")
	}

}

/**
	发送短信验证码
 */
func (s *Sms) Send() SmsResult {
	used_sms_number := int64(0) // 已消耗的短信条数
	countContentCharge()        // 每发一个人需要扣除的短信条数
	mobile_arr := strings.Split(smsapi.GetAddress(), ",")
	sendList := []string{}
	count := len(mobile_arr)
	for k, v := range mobile_arr { // 循环插入到队列，当满足发送数量时再批量发送
		if len(v) > 0 {
			sendList = append(sendList, v)
			// 如果到了指定的数量 或 到达队列的尾部则开始发送
			if int64(len(sendList)) > batSendNum || (k+1) == count {
				// 短信发送
				smsapi.SetAddress(strings.Join(sendList, ","))
				res := smsapi.Send()
				if false == res.Status {
					// 发送失败写入短信发送日志
					//将报错信息写入数据库，可查库排错
					_, err := db.Session().Update("sms_log").SetMap(map[string]interface{}{
						"memo":   res.Info,
						"status": "2",
						"used":   used_sms_number,
					}).Where("smsid = ?", Smsid).Exec()
					if nil != err {
						beego.Error("Update(sms_log) err = ", err, " smsid = ", Smsid)
					}
					res.Info = "发送失败消耗短信条数:" + string(used_sms_number)
					res.Used = used_sms_number
					return res
				}
				used_sms_number = int64(len(sendList)) + 1
				sendList = []string{}
			}
		}
	}

	// 发送成功写入短信发送日志
	_, err := db.Session().Update("sms_log").SetMap(map[string]interface{}{
		"memo":   "发送成功",
		"status": "1",
		"used":   used_sms_number,
	}).Where("smsid = ?", Smsid).Exec()
	if nil != err {
		beego.Error("Update(sms_log) err = ", err, " smsid = ", Smsid)
	}
	return SmsResult{Status: true, Info: "发送成功消耗短信条数:" + string(used_sms_number), Used: used_sms_number}
}

/*
* 计算当前发送内容每次需要扣除的收费条数
* return int
*/
func countContentCharge() float64 {

	number_step_reduce := float64(1)                     //每发一个人需要扣除的短信条数
	content_charge := float64(smsapi.GetContentCharge()) //单封短信最大多少字
	content_size := float64(len(smsapi.GetContent()))    //当前短信有多少字
	if content_charge < content_size {
		number_step_reduce = math.Ceil(content_size / content_size)
	}
	return number_step_reduce
}

func init() {
	smsMap = make(map[string]interface{})
	smsMap["M5cSms"] = new(M5cSms)
}

//产生新编号
func createSmsId() string {
	date := beego.Date(time.Now(), "YmdHis")
	return date + strconv.FormatInt(tools.Random(100000, 999999), 10)
}

//判断是否为合法的短信类型
func checkSmsName(smsname string) bool {
	if _, ok := smsMap[smsname]; ok {
		return true
	} else {
		return false
	}
}

func printMessage(msg string) {
	if "9999" == msg {

	} else if "0000" == msg {

	} else {

	}
	//if($msg == '9999') {
	//	echo 'RespCode=9999|JumpURL=';
	//	exit;
	//} elseif($msg == '0000') {
	//	echo 'RespCode=0000|JumpURL=';
	//	exit;
	//} else {
	//	echo '<!DOCTYPE html><html xmlns="http://www.w3.org/1999/xhtml"><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /></head><body>' . $msg . '</body></html>';
	//}
}

//创建发送报告
func createReport(sendUsername string, smsname string) {

	log := models.SmsLog{}
	data := map[string]interface{}{}
	data["content"] = Content
	data["sms_class"] = smsname
	data["sms_name"] = smsapi.GetName()
	data["w_time"] = time.Now().Unix()
	data["username"] = sendUsername
	data["addressee"] = Address
	data["status"] = 0
	data["memo"] = ""
	//根据smsid查询记录
	err := db.Session().Select("*").From("sms_log").Where("smsid = ?", Smsid).LoadOne(&log)
	if nil != err {
		//没查到，新增一条
		data["smsid"] = Smsid
		_, err := db.Session().InsertInto("sms_log").Map(data).Exec()
		if nil != err {
			beego.Error("InsertInto(sms_log) error = ", err, " smsid = ", Smsid)
		}
	} else {
		//查到，修改
		_, err := db.Session().Update("sms_log").Where("smsid = ?", Smsid).SetMap(data).Exec()
		if nil != err {
			beego.Error("Update(sms_log) error = ", err, " smsid = ", Smsid)
		}
	}
}
