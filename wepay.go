package anypay

import (
	"github.com/panghu1024/anypay/tools"
)

//接口MAP
var ApiUrlMap  = map[string]string{
	"UnifiedOrder":		"https://api.mch.weixin.qq.com/pay/unifiedorder",
	"Refund":			"https://api.mch.weixin.qq.com/secapi/pay/refund",
	"OrderQuery":		"https://api.mch.weixin.qq.com/pay/orderquery",
	"RefundQuery":		"https://api.mch.weixin.qq.com/pay/refundquery",
	"CloseOrder":		"https://api.mch.weixin.qq.com/pay/closeorder",
	"TransferBalance":	"https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers",
	"TransferBank":		"https://api.mch.weixin.qq.com/mmpaysptrans/pay_bank",
}

//微信支付配置
type WeConfig struct {
	Key 	string		//商户KEY
	AppId 	string		//APPID
	MchId 	string		//商户号
	CertKeyPath string  //证书key路径
	CertP12Path string	//证书p12路径
	CertPemPath string	//证书pem路径
}


//微信支付结构体
type WePay struct {
	config WeConfig
}

//创建支付实例
func NewWePay(config WeConfig) (WePay){

	if config.AppId == ""{
		panic("Appid can not be nil")
	}

	if config.Key == ""{
		panic("Key can not be nil")
	}

	if config.MchId == ""{
		panic("MchId can not be nil")
	}

	return WePay{config:config}
}

//签名检查
func (w WePay) SignCheck(data map[string]interface{}) ReturnParam {

	var signStr string
	if _, ok := data["sign"]; ok {
		signStr = data["sign"].(string)
	}else{
		return ReturnParam{-1,"签名参数不存在",nil}
	}

	sign,err := tools.GenerateSignString(data,w.config.Key)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	if sign == signStr{
		return ReturnParam{1,"ok",nil}
	}else{
		return ReturnParam{-1,"验证失败",nil}
	}
}