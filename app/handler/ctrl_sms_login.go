package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"project/app/handler/pkg/e"
	"project/app/service"
)

// LoginSmsCtrl 发送登录短信验证码
type LoginSmsCtrl struct {
	smsService service.ISms
}

func NewLoginSmsCtrl(smsService service.ISms) *LoginSmsCtrl {
	return &LoginSmsCtrl{smsService: smsService}
}

// Send 发送登录短信验证码
func (ctrl *LoginSmsCtrl) Send(c *gin.Context) {
	// 参数绑定与校验
	type Form struct {
		CnCellPhoneNumber string `form:"cn_cell_phone_number" json:"cn_cell_phone_number" binding:"required,cnCellPhoneNumber"`
	}
	var form Form
	if !mustBind(c, &form) {
		return
	}

	// 发送登录短信验证码
	if err := ctrl.smsService.Send(form.CnCellPhoneNumber); err != nil {
		if err, ok := err.(*service.SmsRequestOutOfLimitError); ok {
			fail(c, err, e.CodeResourceExhausted,
				&e.QuotaFailure{Violations: []*e.QuotaFailureViolation{
					&e.QuotaFailureViolation{
						Description: err.Message,
					},
				}},
				&e.RetryInfo{RetryDelay: err.RetryDelay},
			)
		}
		fail(c, errors.Wrap(err, "发送登录验证码失败"), e.CodeInternal)
		return
	} else {
		success(c, nil)
	}
}
