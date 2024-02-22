package email

import (
	"fmt"
	"testing"
)

func TestSendEmail(t *testing.T) {
	url := "http://www.baidu.com"
	body := `<a href="` + url + `"> 点我继续 </a><br><br>若以上链接无法点击，请复制下面内容到浏览器地址栏并访问<br>` + url
	err := SendEmail("wangcunlu2010@163.com", "测试123", body)
	fmt.Println(err)
}

