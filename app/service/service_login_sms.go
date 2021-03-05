package service

import (
	"github.com/patrickmn/go-cache"
	"project/app/pkg/sms"
	"project/app/pkg/util"
	"time"
)

type LoginSmsService struct {
	sender sms.Sender
	cache  *cache.Cache
}

var _ ISms = new(LoginSmsService)

func NewLoginSmsService(sender sms.Sender, cache *cache.Cache) *LoginSmsService {
	return &LoginSmsService{
		sender: sender,
		cache:  cache,
	}
}

const (
	loginSmsKeyPrefix = "login_cell_phone_number:"
	loginSmsExpire    = 5
)

func (service *LoginSmsService) Send(cnCellPhoneNumber string) error {
	ok, retryDelay := service.sendSpeedLimit(cnCellPhoneNumber)
	if !ok {
		return &SmsRequestOutOfLimitError{
			Message:    "登录短信验证码发送频率超限，5 分钟内最多发送 10 次",
			RetryDelay: retryDelay,
		}
	}

	code := util.GenerateRandomDigits(6)
	if err := service.sender.Send(cnCellPhoneNumber, code, loginSmsExpire); err != nil {
		return nil
	}
	service.cache.Set(loginSmsKeyPrefix+cnCellPhoneNumber, cnCellPhoneNumber, loginSmsExpire)
	return nil
}

// 检测客户端请求登录短信验证码发送接口频率是否超限
func (service *LoginSmsService) sendSpeedLimit(cnCellPhoneNumber string) (ok bool, retryDelay time.Duration) {
	// @todo 以下为 mock 代码
	boolean := time.Now().UnixNano()%2 == 0
	if boolean {
		return true, 0
	} else {
		return false, time.Minute * 5
	}
}

// 验证登录短信验证码是否有效
func (service *LoginSmsService) Verify(cnCellPhoneNumber string, code string) bool {
	data, ok := service.cache.Get(loginSmsKeyPrefix + cnCellPhoneNumber)
	if ok && data.(string) == code {
		return true
	} else {
		return false
	}
}
