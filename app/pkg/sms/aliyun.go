package sms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AliyunLoginSms 负责发送“登录”短信验证码（使用 aliyun sms）
type AliyunLoginSms struct {
	accessKeyId     string
	accessKeySecret string
	regionId        string
	signName        string
	templateCode    string
}

var _ Sender = new(AliyunLoginSms)

func NewAliyunLoginSms(v *viper.Viper) *AliyunLoginSms {
	return &AliyunLoginSms{
		accessKeyId:     v.GetString("aliyunLoginSms.accessKeyId"),
		accessKeySecret: v.GetString("aliyunLoginSms.accessKeySecret"),
		regionId:        v.GetString("aliyunLoginSms.regionId"),
		signName:        v.GetString("aliyunLoginSms.signName"),
		templateCode:    v.GetString("aliyunLoginSms.templateCode"),
	}
}

// Send 发送登录短信验证码
func (sms *AliyunLoginSms) Send(cellPhoneNumber string, code string, expire int) error {
	client, err := dysmsapi.NewClientWithAccessKey(sms.regionId, sms.accessKeyId, sms.accessKeySecret)
	if err != nil {
		return errors.Wrap(err, "AliyunLoginSms new client failed")
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = sms.signName
	request.TemplateCode = sms.templateCode
	request.TemplateParam = `{"code":"` + code + `"}` // 这里也可以把 expire 配置进模板
	request.PhoneNumbers = cellPhoneNumber
	_, err = client.SendSms(request)
	if err != nil {
		return errors.Wrap(err, "AliyunLoginSms client send sms failed")
	}

	return nil
}
