package verifycode

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"gohub/pkg/sms"
	"strings"
	"sync"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once

var internalVerifyCode *VerifyCode

func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				KeyPrefix:   config.GetString("app.name") + ":verifycode:",
			},
		}
	})
	return internalVerifyCode
}

func (vc *VerifyCode) SendSMS(phone string) bool {
	code := vc.generateVerifyCode(phone)
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.sms_debug_phone_prefix")) {
		return true
	}

	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.tencent.template_code"),
		Data:     map[string]string{"1": code, "2": config.GetString("verifycode.code_expire_minutes")},
	})
}

func (vs *VerifyCode) CheckAnswer(key, answer string) bool {
	logger.DebugJSON("验证码", "验证验证码", map[string]string{key: answer})

	if !app.IsProduction() && strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix")) || strings.HasPrefix(key, config.GetString("verifycode.debug_email_suffix")) {
		return answer == config.GetString("verifycode.debug_code")
	}

	return vs.Store.Verify(key, answer, true)
}

func (vs *VerifyCode) generateVerifyCode(key string) string {
	code := helpers.RandonNumber(config.GetInt("verifycode.code_length"))

	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})

	vs.Store.Set(key, code)

	return code
}
