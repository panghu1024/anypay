package anypay

import (
	"io/ioutil"
	"encoding/xml"
	"github.com/panghu1024/anypay/tools"
	"net/http"
	"strings"
)

//关闭订单参数
type WeCloseOrderParam struct {
	Appid string `xml:"appid"`					// appid
	MchId string `xml:"mch_id"`					// 商户id
	NonceStr string `xml:"nonce_str"`			// 随机字符串
	Sign string `xml:"sign"`					// 签名
	SignType string `xml:"sign_type"`			// 加密方式
	OutTradeNo string `xml:"out_trade_no"`		// 商户订单号
}


//关闭订单返回结构体
type WeResCloseOrder struct {
	ReturnCode string `xml:"return_code"`				// SUCCESS/FAIL,此字段是通信标识，表示接口层的请求结果，并非退款状态
	ReturnMsg string `xml:"return_msg"`					// 当return_code为FAIL时返回信息为错误原因
	ResultCode string `xml:"result_code"`				// SUCCESS/FAIL，SUCCESS退款申请接收成功，结果通过退款查询接口查询,FAIL 提交业务失败
	ErrCode string `xml:"err_code"`						// 列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"`				// 结果信息描述
	Appid string `xml:"appid"`							// 微信分配的公众账号ID
	MchId string `xml:"mch_id"`							// 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"`					// 随机字符串，不长于32位
	Sign string `xml:"sign"`							// 签名
	OutTradeNo string `xml:"out_trade_no"`				// 商户订单号
}


//关闭订单
func (w WePay) CloseOrder(closeParam WeCloseOrderParam) ReturnParam {
	nonceStr := tools.GenerateNonceString()

	closeParam.Appid = w.config.AppId	//设置APP ID
	closeParam.NonceStr = nonceStr		//设置随机字符串
	closeParam.MchId = w.config.MchId	//设置商户ID

	queryMap := tools.Struct2Map(closeParam)

	sign,err := tools.GenerateSignString(queryMap,w.config.Key)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	closeParam.Sign = sign

	requestXml,err := tools.GenerateRequestXml(closeParam)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	//发起请求
	r,err := http.Post(ApiUrlMap["CloseOrder"],"text/xml",strings.NewReader(requestXml))

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	var res WeResCloseOrder

	xml.Unmarshal([]byte(string(body)),&res)

	if res.ReturnCode != "SUCCESS" || res.ResultCode != "SUCCESS"{
		return ReturnParam{-1,res.ReturnMsg,res}
	}

	return ReturnParam{1,"success",res}
}