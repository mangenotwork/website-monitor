package dao

import (
	"fmt"
	"strings"

	"website-monitor/master/entity"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/mail"
	"github.com/mangenotwork/common/utils"
)

type MailEr interface {
	// SetMail 设置邮件
	SetMail(data *entity.Mail) error

	// GetMail 获取邮件
	GetMail() (*entity.Mail, error)

	// IsMail 是否配置邮件
	IsMail() bool

	// Send 发送邮件
	Send(title, body string)

	// Check 检查实例数据
	Check(data *entity.Mail) error
}

func NewMail() MailEr {
	return new(mailDao)
}

type mailDao struct{}

func (m *mailDao) SetMail(data *entity.Mail) error {
	return DB.Set(MailTable, MailConfKeyName, data)
}

func (m *mailDao) GetMail() (*entity.Mail, error) {
	data := &entity.Mail{}
	err := DB.Get(MailTable, MailConfKeyName, data)
	return data, err
}

func (m *mailDao) IsMail() bool {
	data, err := m.GetMail()
	if err != nil {
		log.Error(err)
	}
	if len(data.Host) > 0 && len(data.From) > 0 && len(data.AuthCode) > 0 {
		return true
	}
	return false
}

func (m *mailDao) Send(title, body string) {
	mailInfo, _ := m.GetMail()
	if !mailInfo.Open { // 关闭
		return
	}

	data, err := m.GetMail()
	if err != nil {
		log.Error(err)
		return
	}

	send := mail.NewMail(data.Host, data.From, data.AuthCode, data.Port)
	toList := strings.Split(data.ToList, ";")

	err = send.Title(title).HtmlBody(body).SendMore(toList)
	if err != nil {
		log.Error(err)
		return
	}

}

func (m *mailDao) Check(data *entity.Mail) error {
	errFormat := "邮件设置参数错误: %s"

	if len(data.From) < 1 {
		return fmt.Errorf(errFormat, "发件人不能为空!")
	}

	if len(data.Host) < 1 {
		return fmt.Errorf(errFormat, "邮件服务器不能为空!")
	}

	if len(data.AuthCode) < 1 {
		return fmt.Errorf(errFormat, "邮件服务授权码不能为空!")
	}

	if len(data.ToList) < 1 {
		return fmt.Errorf(errFormat, "通知收件人不能为空!")
	}

	if data.Port == 0 {
		data.Port = 25
	}

	data.ToList = utils.CleaningStr(data.ToList)
	data.ToList = strings.Replace(data.ToList, "；", ";", -1)

	return nil
}
