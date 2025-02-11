package sms

import (
	"gohub/pkg/logger"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Tencent struct{}

func (s *Tencent) Send(phone string, message Message, config map[string]string) bool {
	// 硬编码密钥到代码中有可能随代码泄露而暴露，有安全隐患，并不推荐。
	// 为了保护密钥安全，建议将密钥设置在环境变量中或者配置文件中，请参考本文凭证管理章节。

	logger.DebugJSON("短信[腾讯云]", "配置信息", config)

	// 1. 可以使用 NewCredential 来创建一个普通的密钥
	credential := common.NewCredential(
		config["access_key_id"],
		config["access_key_secret"],
	)

	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()
	/* SDK默认使用POST方法。
	 * 如果您一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 10 // 请求超时时间，单位为秒(默认60秒)
	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"
	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考 https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr("1400549759")
	request.SignName = common.StringPtr(config["sign_name"])
	request.TemplateId = common.StringPtr(message.Template) //1040656
	paramSet := []string{}
	for _, v := range message.Data {
		paramSet = append(paramSet, v)
	}
	request.TemplateParamSet = common.StringPtrs(paramSet)
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})

	response, err := client.SendSms(request)
	// 处理异常

	logger.DebugJSON("短信[腾讯云]", "请求内容", request)
	logger.DebugJSON("短信[腾讯云]", "接口响应", response)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.ErrorString("短信[腾讯云]", "发信失败", err.Error())
		return false
	}
	if err != nil {
		logger.ErrorString("短信[腾讯云]", "解析响应 JSON 错误", err.Error())
	}
	logger.DebugString("短信[腾讯云]", "发信成功", response.ToJsonString())

	return true
}
