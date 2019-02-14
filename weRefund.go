package anypay

import (
	"io/ioutil"
	"encoding/xml"
	"github.com/panghu1024/anypay/tools"
	"crypto/x509"
	"crypto/tls"
	"net/http"
	"bytes"
)

//退款参数
type WeRefundParam struct {
	Appid string `xml:"appid"`					// appid
	MchId string `xml:"mch_id"`					// 商户id
	NonceStr string `xml:"nonce_str"`			// 随机字符串
	Sign string `xml:"sign"`					// 签名
	TransactionId string `xml:"transaction_id"`	// 交易流水号
	OutRefundNo string `xml:"out_refund_no"`	// 退款订单号
	TotalFee string `xml:"total_fee"`			// 总金额
	RefundFee string `xml:"refund_fee"`			// 退款金额
}

//退款返回结构体
type WeResRefund struct {
	ReturnCode string `xml:"return_code"`				// SUCCESS/FAIL,此字段是通信标识，表示接口层的请求结果，并非退款状态
	ReturnMsg string `xml:"return_msg"`					// 当return_code为FAIL时返回信息为错误原因
	ResultCode string `xml:"result_code"`				// SUCCESS/FAIL，SUCCESS退款申请接收成功，结果通过退款查询接口查询,FAIL 提交业务失败
	ErrCode string `xml:"err_code"`						// 列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"`				// 结果信息描述
	Appid string `xml:"appid"`							// 微信分配的公众账号ID
	MchId string `xml:"mch_id"`							// 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"`					// 随机字符串，不长于32位
	Sign string `xml:"sign"`							// 签名
	TransactionId string `xml:"transaction_id"`			// 微信订单号
	OutTradeNo string `xml:"out_trade_no"`				// 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	OutRefundNo string `xml:"out_refund_no"`			// 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔
	RefundId string `xml:"refund_id"`					// 微信退款单号
	RefundChannel string `xml:"refund_channel"`			// 退款渠道
	RefundFee int `xml:"refund_fee"`					// 退款金额单位分
	CouponRefundFee int `xml:"coupon_refund_fee"`		// 优惠券退款金额 单位分
	TotalFee int `xml:"total_fee"`						// 订单总金额 单位分
	CashFee int `xml:"cash_fee"`						// 现金金额 单位分
	CouponRefundCount int `xml:"coupon_refund_count"`	// 优惠券退款数量
	CashRefundFee int `xml:"cash_refund_fee"`			// 现金退款金额 单位分
}

//订单退款
func (w WePay) Refund(refundParam WeRefundParam) ReturnParam {
	nonceStr := tools.GenerateNonceString()

	refundParam.Appid = w.config.AppId	//设置APP ID
	refundParam.NonceStr = nonceStr		//设置随机字符串
	refundParam.MchId = w.config.MchId	//设置商户ID

	refundMap := tools.Struct2Map(refundParam)

	sign,err := tools.GenerateSignString(refundMap,w.config.Key)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	refundParam.Sign = sign

	requestXml,err := tools.GenerateRequestXml(refundParam)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	pool := x509.NewCertPool()

	pemData,err := ioutil.ReadFile(w.config.CertP12Path)

	if err != nil{
		return ReturnParam{-1,"无法读取apiclient_cert.p12文件",nil}
	}

	pool.AppendCertsFromPEM(pemData)

	cert, err := tls.LoadX509KeyPair(w.config.CertPemPath,w.config.CertKeyPath)

	if err != nil{
		return ReturnParam{-1,"无法读取apiclient_cert.pem或apiclient_key.pem文件",nil}
	}

	mTLSConfig:=& tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:pool,
		InsecureSkipVerify: true,
	}

	tr :=& http.Transport{
		TLSClientConfig:mTLSConfig,
	}

	client := & http.Client{
		Transport:tr,
	}

	req, err := http.NewRequest("POST",ApiUrlMap["Refund"],bytes.NewBuffer([]byte(requestXml)))

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	req.Header.Set("Content-Type", "text/xml")
	r,err := client.Do(req)

	defer r.Body.Close()

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	body, err := ioutil.ReadAll(r.Body)

	//
	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	var res WeResRefund

	xml.Unmarshal([]byte(string(body)),&res)

	if res.ReturnCode != "SUCCESS" || res.ResultCode != "SUCCESS"{
		return ReturnParam{-1,res.ReturnMsg,res}
	}

	return ReturnParam{1,"ok",res}
}
