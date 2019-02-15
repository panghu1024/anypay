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

//转账到余额参数
type WeTransferBalanceParam struct {
	MchAppid string `xml:"mch_appid"`			// 申请商户号的AppId
	Mchid string `xml:"mchid"`					// 商户id
	NonceStr string `xml:"nonce_str"`			// 随机字符串
	Sign string `xml:"sign"`					// 签名
	DeviceInfo string `xml:"device_info"`		// 设备信息
	PartnerTradeNo string `xml:"partner_trade_no"`	// 商户单号 需唯一 只能是字母或者数字
	Openid string `xml:"openid"`				// 用户OpenId
	CheckName string `xml:"check_name"`			// NO_CHECK：不校验真实姓名 FORCE_CHECK：强校验真实姓名
	ReUserName string `xml:"re_user_name"`		// 用户真实姓名 如果check_name设置为FORCE_CHECK，则必填用户真实姓名
	Amount string `xml:"amount"`				// 企业付款金额 单位 分
	Desc string `xml:"desc"`					// 付款备注
	SpbillCreateIp string `xml:"spbill_create_ip"`	// 用户端真实IP
}

//转账到余额返回结构体
type WeResTransferBalance struct {
	ReturnCode string `xml:"return_code"`				// SUCCESS/FAIL,此字段是通信标识，表示接口层的请求结果，并非退款状态
	ReturnMsg string `xml:"return_msg"`					// 当return_code为FAIL时返回信息为错误原因
	ResultCode string `xml:"result_code"`				// SUCCESS/FAIL，SUCCESS退款申请接收成功，结果通过退款查询接口查询,FAIL 提交业务失败
	ErrCode string `xml:"err_code"`						// 列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"`				// 结果信息描述
	MchAppid string `xml:"mch_appid"`					// 申请商户的AppId
	MchId string `xml:"mch_id"`							// 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"`					// 随机字符串，不长于32位
	DeviceInfo string `xml:"device_info"`				// 设备信息
	PartnerTradeNo string `xml:"partner_trade_no"`		// 商户订单号，需保持历史全局唯一性
	PaymentNo string `xml:"payment_no"`					// 企业付款成功，返回的微信付款单号
	PaymentTime string `xml:"payment_time"`				// 企业付款成功时间 yyyy-mm-dd hh:ii:ss

}

//参数检查
func (w WePay) paramCheckFormTransferBalance(param WeTransferBalanceParam) bool {
	if param.Openid == ""{
		return false
	}
	if param.PartnerTradeNo == ""{
		return false
	}
	if param.CheckName != "FORCE_CHECK" && param.CheckName != "NO_CHECK"{
		return false
	}
	if param.CheckName == "FORCE_CHECK" && param.ReUserName == ""{
		return false
	}
	if param.Desc == ""{
		return false
	}
	if param.SpbillCreateIp == ""{
		return false
	}

	return true
}

//转账到余额
func (w WePay) TransferBalance(transParam WeTransferBalanceParam) ReturnParam {

	if w.paramCheckFormTransferBalance(transParam) == false{
		return ReturnParam{-1,"请检查必传参数",nil}
	}

	nonceStr := tools.GenerateNonceString()

	transParam.MchAppid = w.config.AppId	//设置APP ID
	transParam.NonceStr = nonceStr		//设置随机字符串
	transParam.Mchid = w.config.MchId	//设置商户ID

	transMap := tools.Struct2Map(transParam)

	sign,err := tools.GenerateSignString(transMap,w.config.Key)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	transParam.Sign = sign

	requestXml,err := tools.GenerateRequestXml(transParam)

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

	req, err := http.NewRequest("POST",ApiUrlMap["TransferBalance"],bytes.NewBuffer([]byte(requestXml)))

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

	var res WeResTransferBalance

	xml.Unmarshal([]byte(string(body)),&res)

	if res.ReturnCode != "SUCCESS" || res.ResultCode != "SUCCESS"{
		return ReturnParam{-1,res.ReturnMsg,res}
	}

	return ReturnParam{1,"ok",res}
}
