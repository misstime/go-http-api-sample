package sms

type Sender interface {
	// 向单个手机号发送短信验证码
	//
	// 参数 cellPhoneNumber 为手机号码，必须支持国外手机号
	// 参数 code 为验证码内容，例：`098909`、`aLNs89`、`9089`
	// 参数 expire 展示给用户的验证码有效期，单位：分钟
	//
	// Note：Send 方法只负责将参数 code、expire 解析到模板中然后返回给用户，其他如
	// 缓存验证码、验证验证码是否正确等操作一律不要出现在该方法内。
	Send(cellPhoneNumber string, code string, expire int) error
}
