package email

import (
	"gopkg.in/gomail.v2"
)

//发送邮件
func SendEmail(to string, subject string, body string) error {
	smtp_host := "email-smtp.us-west-2.amazonaws.com"           //SMTP服务器地址，类似 smtp.tom.com
	smtp_port := 465                                            //SMTP服务器端口
	smtp_user := "AKIAIW7ODO6DOMLIF63A"                         //Smtp认证的用户名
	smtp_pass := "AuicnrYl4TnQzg5pdQLEEiqLgrLdDSP47cX5rVZiBKSP" //Smtp认证的密码，一般等同pop3密码
	from := "support@m.cc"                                      //发信人Email地址，你的发信信箱地址

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)           //发送给用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(smtp_host, smtp_port, smtp_user, smtp_pass)
	err := d.DialAndSend(m)
	return err
}
