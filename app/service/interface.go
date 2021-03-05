package service

import "time"

// 客户端短信验证码发送请求频率超限错误
type SmsRequestOutOfLimitError struct {
	// 错误描述
	Message string
	// 客户端下一次请求发送短信验证码的最短时间间隔
	RetryDelay time.Duration
}

func (err *SmsRequestOutOfLimitError) Error() string {
	return err.Message
}

// 短信验证码服务接口
type ISms interface {
	// 向单个手机号发送短信验证码
	//
	// 客户端请求频率超限时，返回 *SmsRequestOutOfLimitError
	Send(cellPhoneNumber string) error
	// 验证短信验证码是否有效
	Verify(cellPhoneNumber, code string) bool
}
