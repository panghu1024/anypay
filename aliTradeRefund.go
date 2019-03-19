package anypay

import (
	"net/url"
	"errors"
	"github.com/panghu1024/anypay/tools"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

//退款参数
type AliTradeRefundParam struct {
	// - - - Base Param
	AppId string `json:"app_id"`								//支付宝分配给开发者的应用ID
	Method string `json:"method"`								//接口名称
	Format string `json:"format"`								//仅支持JSON
	Charset string `json:"charset"`								//请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType string `json:"sign_type"`							//商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign string `json:"sign"`									//商户请求参数的签名串
	Timestamp string `json:"timestamp"`							//发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
	Version string `json:"version"`								//调用的接口版本，固定为：1.0
	AppAuthToken string `json:"app_auth_token"`					//详见应用授权概述
	BizContent TradeRefundBizContent `json:"biz_content"`						//请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
}

//请求参数
type TradeRefundBizContent struct {
	OutTradeNo string `json:"out_trade_no"`						//订单支付时传入的商户订单号,不能和 trade_no同时为空。
	TradeNo string `json:"trade_no"`							//支付宝交易号，和商户订单号不能同时为空
	RefundAmount string `json:"refund_amount"`					//需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundCurrency string `json:"refund_currency"`				//订单退款币种信息
	RefundReason string `json:"refund_reason"`					//退款的原因说明
	OutRequestNo string `json:"out_request_no"`					//标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorId string `json:"operator_id"`						//商户的操作员编号
	StoreId string `json:"store_id"`							//商户的门店编号
	TerminalId string `json:"terminal_id"`						//商户的终端编号
	GoodsDetail []GoodsDetail `json:"goods_detail"`				//退款包含的商品列表信息，Json格式。 其它说明详见：“商品明细说明”
	RefundRoyaltyParameters []RefundRoyaltyParameters			//退分账明细信息
	OrgPid string												//银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要退款的交易所属收单机构的pid;
}

//退款包含的商品列表信息
type GoodsDetail struct {
	GoodsId string `json:"goods_id"`				//商品的编号
	AlipayGoodsId string `json:"alipay_goods_id"`	//支付宝定义的统一商品编号
	GoodsName string `json:"goods_name"`			//商品名称
	Quantity string `json:"quantity"`				//商品数量
	Price string `json:"price"`						//商品单价，单位为元
	GoodsCategory string `json:"goods_category"`	//商品类目
	CategoriesTree string `json:"categories_tree"`	//商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割
	Body string `json:"body"`						//商品描述信息
	ShowUrl string `json:"show_url"`				//商品的展示地址
}

//退分账明细信息
type RefundRoyaltyParameters struct {
	RoyaltyType string `json:"royalty_type"`			//分账类型. 普通分账为：transfer;	补差为：replenish;	为空默认为分账transfer;
	TransOut string `json:"trans_out"`					//支出方账户。如果支出方账户类型为userId，本参数为支出方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果支出方类型为loginName，本参数为支出方的支付宝登录号；
	TransOutType string `json:"trans_out_type"`			//支出方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;loginName表示是支付宝登录号；
	TransInType string `json:"trans_in_type"`			//收入方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;cardSerialNo表示是卡编号;loginName表示是支付宝登录号；
	TransIn string `json:"trans_in"`					//收入方账户。如果收入方账户类型为userId，本参数为收入方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果收入方类型为cardSerialNo，本参数为收入方在支付宝绑定的卡编号；如果收入方类型为loginName，本参数为收入方的支付宝登录号；
	Amount float64 `json:"amount"`						//分账的金额，单位为元
	AmountPercentage float64 `json:"amount_percentage"`	//分账信息中分账百分比。取值范围为大于0，少于或等于100的整数。
	Desc string `json:"desc"`							//分账描述
}

//退款返回结构体
type AliResTradeRefund struct {
	AlipayTradeRefundResponse AlipayTradeRefundResponse `json:"alipay_trade_refund_response"`
	Sign string `json:"sign"`														//签名,详见文档
}

//返回参数
type AlipayTradeRefundResponse struct {
	Code string `json:"code"`														//网关返回码,详见文档
	Msg string `json:"msg"`															//网关返回码描述,详见文档
	SubCode string `json:"sub_code"`												//业务返回码，参见具体的API接口文档
	SubMsg string `json:"sub_msg"`													//业务返回码描述，参见具体的API接口文档
	OutTradeNo string `json:"out_trade_no"`											//商户网站唯一订单号
	TradeNo string `json:"trade_no"`												//该交易在支付宝系统中的交易流水号。最长64位。
	BuyerLogonId string `json:"buyer_logon_id"`
	FundChange string `json:"fund_change"`
	RefundFee string `json:"refund_fee"`
	RefundCurrency string `json:"refund_currency"`
	GmtRefundPay string `json:"gmt_refund_pay"`
	RefundDetailItemList TradeFundBill `json:"refund_detail_item_list"`
	StoreName string `json:"store_name"`
	BuyerUserId string `json:"buyer_user_id"`
	RefundPresetPaytoolList PresetPayToolInfo `json:"refund_preset_paytool_list"`
	RefundChargeAmount string `json:"refund_charge_amount"`
	RefundSettlementId string `json:"refund_settlement_id"`
	PresentRefundBuyerAmount string `json:"present_refund_buyer_amount"`
	PresentRefundDiscountAmount string `json:"present_refund_discount_amount"`
	PresentRefundMdiscountAmount string `json:"present_refund_mdiscount_amount"`
}

type TradeFundBill struct {
	FundChannel string
	Amount string
	RealAmount string
	FundType string
}

type PresetPayToolInfo struct {
	Amount [] string
	AssertTypeCode string
}

//退款
func (ali AliPay) TradeRefund(refundParam AliTradeRefundParam) ReturnParam {

	//处理参数
	param,err := ali.tradeRefundParam(refundParam)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	//对参数进行排序
	src,_ := tools.SortData(param)

	var hash crypto.Hash
	if refundParam.SignType == "RSA" {
		hash = crypto.SHA1
	} else {
		hash = crypto.SHA256
	}

	privateKey := tools.ParsePrivateKey(ali.config.PrivateKeyString)

	sign, err := ali.SignPKCS1v15([]byte(src), privateKey, hash)

	signWithBase64 := base64.StdEncoding.EncodeToString(sign)

	param.Add("sign", signWithBase64)

	url := AliCommonApi+"?"+param.Encode()

	r,err := http.Get(url)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	var res AliResTradeRefund

	json.Unmarshal([]byte(string(body)),&res)

	if res.AlipayTradeRefundResponse.Code != "10000"{
		return ReturnParam{-1,res.AlipayTradeRefundResponse.Msg,res}
	}

	return ReturnParam{1,"ok",res}
}

//处理交易参数
func (ali AliPay) tradeRefundParam(orderParam AliTradeRefundParam) (url.Values,error){
	var param = url.Values{}
	param.Add("app_id", ali.config.AppId)
	param.Add("method", "alipay.trade.refund")
	param.Add("format", "JSON")
	param.Add("charset", orderParam.Charset)
	param.Add("sign_type", orderParam.SignType)
	param.Add("timestamp", orderParam.Timestamp)
	param.Add("version", "1.0")

	if orderParam.AppAuthToken != ""{
		param.Add("app_auth_token", orderParam.AppAuthToken)	//
	}

	bizContent:=tools.Struct2Map(orderParam.BizContent)

	if orderParam.BizContent.GoodsDetail != nil{
		bizContent["goods_detail"] = tools.Struct2Map(orderParam.BizContent.GoodsDetail)
	}

	if orderParam.BizContent.RefundRoyaltyParameters != nil{
		bizContent["refund_royalty_parameters"] = tools.Struct2Map(orderParam.BizContent.RefundRoyaltyParameters)
	}

	bizContentString,err := json.Marshal(bizContent)

	if err != nil{
		return nil,errors.New("生成biz_content字符串失败")
	}

	param.Add("biz_content",string(bizContentString))

	return param,nil
}
