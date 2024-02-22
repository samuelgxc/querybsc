package server

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"shangchenggo/tools"
	"time"
)

func StartBeegoTplfunc()  {
	beego.AddFuncMap("i18n", i18n.Tr)     //i18n国际化语言模板函数
	beego.AddFuncMap("md5", tools.Md5)
	beego.AddFuncMap("dateFormatNow", dateFormatNow)
	beego.AddFuncMap("dateFormatUnix", dateFormatUnix)
}

func dateFormatNow(format string) string {
	return time.Now().Format(format)
}

func dateFormatUnix(unix int64, format string) string {
	return tools.FormatTimestamp(unix, format)
}
