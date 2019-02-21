package anypay

import (
	"encoding/xml"
	"io/ioutil"
	"bytes"
	"net/http"
	"github.com/panghu1024/anypay/tools"
)

//通知结构体
type WeResNotify struct {
	ReturnCode string `xml:"return_code"`				// SUCCESS/FAIL,此字段是通信标识，表示接口层的请求结果，并非退款状态
	ReturnMsg string `xml:"return_msg"`					// 当return_code为FAIL时返回信息为错误原因
	ResultCode string `xml:"result_code"`				// SUCCESS/FAIL，SUCCESS退款申请接收成功，结果通过退款查询接口查询,FAIL 提交业务失败
	ErrCode string `xml:"err_code"`						// 列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"`				// 结果信息描述
	Appid string `xml:"appid"`							// 公众号APP ID
	MchId string `xml:"mch_id"`							// 商户ID
	DeviceInfo string `xml:"device_info"`				// 设备信息
	NonceStr string `xml:"nonce_str"`					// 随机字符串
	Sign string `xml:"sign"`							// 签名
	SignType string `xml:"sign_type"`					// 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	Openid string `xml:"openid"`						// 公众号OpenId
	IsSubscribe string `xml:"is_subscribe"`				// 用户是否关注公众账号，Y-关注，N-未关注
	TradeType string `xml:"trade_type"`					// 交易类型 JSAPI、NATIVE、APP等
	BankType string `xml:"bank_type"`					// 银行类型，采用字符串类型的银行标识
	TotalFee int `xml:"total_fee"`						// 订单总金额，单位为分
	FeeType string `xml:"fee_type"`						// 货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY
	CashFee string `xml:"cash_fee"`						// 现金支付金额订单现金支付金额
	CashFeeType string `xml:"cash_fee_type"`			// 货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY
	CouponFee int `xml:"coupon_fee"`					// 代金券金额<=订单金额，订单金额-代金券金额=现金支付金额
	CouponCount int `xml:"coupon_count"`				// 代金券使用数量
	TransactionId string `xml:"transaction_id"`			// 微信支付订单号
	OutTradeNo string `xml:"out_trade_no"`				// 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	Attach string `xml:"attach"`						// 商家数据包，原样返回
	TimeEnd string `xml:"time_end"`						// 支付完成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010
}

//获取通知数据
func (w WePay) Notify(r *http.Request,checkAppId bool) ReturnParam {

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var resNotify WeResNotify

	xml.Unmarshal(body, &resNotify)

	if checkAppId == true && resNotify.Appid != w.config.AppId{
		return ReturnParam{-1,"AppId不匹配!",resNotify}
	}

	if resNotify.MchId != w.config.MchId{
		return ReturnParam{-1,"商户Id不匹配!",resNotify}
	}

	resMap := tools.Struct2Map(resNotify)

	signCheck := w.SignCheck(resMap)

	if signCheck.Status != 1{
		return ReturnParam{-1,"签名验签未通过",resMap}
	}

	return ReturnParam{1,"ok",resNotify}
}