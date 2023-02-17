package utils

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	//EMAIL_USER = "2682435647@qq.com"
	EMAIL_PWD = "vhithufhtnqkdjac"
	EMAIL_HOST = "smtp.qq.com"
	EMAIL_PORT = "25"
)
//emailUser：企业邮箱
//from：发件人邮箱，格式为“名称+<邮箱>”，也可以直接写邮箱，默认显示的发件人为@符号前的名称
//to：收件人邮箱地址
//subject：邮件标题
//text：邮件正文
func SendEmail(xlsxFileName, emailUser, from, subject,text string,toEmail []string) error{
	e := email.NewEmail()

	//CC全称是Carbon Copy，意为抄送，BCC全称Blind Carbon Copy，意为暗抄送，收件人看不到被暗抄送给了谁。
	//e.Cc = []string{"xxxxxxx@qq.com"}
	//e.Bcc = []string{"xxxxxxx@qq.com"}

	//e.From = "小礼鱼名品<"+EMAIL_USER+">"
	e.From = from +"<"+ emailUser +">"
	e.To = toEmail
	e.Subject = subject
	e.Text = []byte(text)

	//发送附件非常简单，直接传入文件名即可
	if xlsxFileName != ""{
		_,err := e.AttachFile(xlsxFileName)
		if err != nil{
			return err
		}
	}

	//参数1：通常，identity应该是空字符串，以用作用户名。
	//参数2：用户名
	//参数3：密码，如果拿到了授权码，则填写授权码
	//参数4：服务器地址，163的地址是smtp.163.com，其他平台可自行查看
	err := e.Send(EMAIL_HOST+":"+EMAIL_PORT, smtp.PlainAuth("", emailUser, EMAIL_PWD, EMAIL_HOST))
	if err != nil{
		return err
	}
	return nil
}
