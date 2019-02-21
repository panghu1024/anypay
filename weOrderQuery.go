package anypay

import (
	"github.com/panghu1024/anypay/tools"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/xml"
)

//退款参数
type WeOrderQueryParam struct {
	Appid string `xml:"appid"`					// appid
	MchId string `xml:"mch_id"`					// 商户id
	NonceStr string `xml:"nonce_str"`			// 随机字符串
	Sign string `xml:"sign"`					// 签名
	TransactionId string `xml:"transaction_id"`	// 交易流水号
	OutTradeNo string `xml:"out_trade_no"`		// 商户订单号
	SignType string `xml:"sign_type"`			// 签名类型
}

//退款返回结构体
type WeResOrderQuery struct {
	ReturnCode string `xml:"return_code"`				// SUCCESS/FAIL,此字段是通信标识，表示接口层的请求结果，并非退款状态
	ReturnMsg string `xml:"return_msg"`					// 当return_code为FAIL时返回信息为错误原因
	ResultCode string `xml:"result_code"`				// SUCCESS/FAIL，SUCCESS退款申请接收成功，结果通过退款查询接口查询,FAIL 提交业务失败
	ErrCode string `xml:"err_code"`						// 列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"`				// 结果信息描述
	Appid string `xml:"appid"`							// 微信分配的公众账号ID
	MchId string `xml:"mch_id"`							// 微信支付分配的商户号
	DeviceInfo string `xml:"device_info"`				// 设备信息
	NonceStr string `xml:"nonce_str"`					// 随机字符串，不长于32位
	Sign string `xml:"sign"`							// 签名
	Openid string `xml:"openid"`						// 微信openid
	IsSubscribe string `xml:"is_subscribe"`				// 是否关注公众号
	TradeType string `xml:"trade_type"`					// 交易类型
	BankType string `xml:"bank_type"`					// 银行类型
	TransactionId string `xml:"transaction_id"`			// 微信订单号
	OutTradeNo string `xml:"out_trade_no"`				// 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	TotalFee int `xml:"total_fee"`						// 订单总金额 单位分
	FeeType string `xml:"fee_type"`						// 货币类型 CNY
	Attach string `xml:"attach"`						// 附加信息
	TimeEnd string `xml:"time_end"`						// 支付完成时间
	TradeState string `xml:"trade_state"`				// 交易状态 SUCCESS—支付成功,REFUND—转入退款,NOTPAY—未支付,CLOSED—已关闭,REVOKED—已撤销,USERPAYING--用户支付中,PAYERROR--支付失败
}

//订单查询
func (w WePay) OrderQuery(queryParam WeOrderQueryParam) ReturnParam {

	nonceStr := tools.GenerateNonceString()

	queryParam.Appid = w.config.AppId	//设置APP ID
	queryParam.NonceStr = nonceStr		//设置随机字符串
	queryParam.MchId = w.config.MchId	//设置商户ID

	queryMap := tools.Struct2Map(queryParam)

	sign,err := tools.GenerateSignString(queryMap,w.config.Key)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	queryParam.Sign = sign

	requestXml,err := tools.GenerateRequestXml(queryParam)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	//发起请求
	r,err := http.Post(ApiUrlMap["OrderQuery"],"text/xml",strings.NewReader(requestXml))

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	var res WeResOrderQuery

	xml.Unmarshal([]byte(string(body)),&res)

	if res.ReturnCode != "SUCCESS" || res.ResultCode != "SUCCESS"{
		return ReturnParam{-1,res.ReturnMsg,res}
	}

	return ReturnParam{1,"success",res}

}