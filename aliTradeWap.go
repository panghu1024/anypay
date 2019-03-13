package anypay

import (
	"github.com/panghu1024/anypay/tools"
	"net/url"
	"errors"
	"crypto"
	"encoding/base64"
	"encoding/json"
)

//下单参数
type AliTradeWapParam struct {
	AppId string `json:"app_id"`								//支付宝分配给开发者的应用ID
	Method string `json:"method"`								//接口名称
	Format string `json:"format"`								//仅支持JSON
	ReturnUrl string `json:"return_url"`						//HTTP/HTTPS开头字符串
	Charset string `json:"charset"`								//请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType string `json:"sign_type"`							//商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign string `json:"sign"`									//商户请求参数的签名串
	Timestamp string `json:"timestamp"`							//发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
	Version string `json:"version"`								//调用的接口版本，固定为：1.0
	NotifyUrl string `json:"notify_url"`						//支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	AppAuthToken string `json:"app_auth_token"`					//详见应用授权概述
	BizContent TradeWapBizContent `json:"biz_content"`			//请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
}

type TradeWapBizContent struct {
	Body string `json:"body"`									//商品描述
	Subject string `json:"subject"`								//商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo string `json:"out_trade_no"`						//商户网站唯一订单号
	TimeoutExpress string `json:"timeout_express"`				//该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	TimeExpire string `json:"time_expire"`						//绝对超时时间，格式为yyyy-MM-dd HH:mm。
	TotalAmount string `json:"total_amount"`					//订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	AuthToken string `json:"auth_token"`						//针对用户授权接口，获取用户相关数据时，用于标识用户授权关系
	GoodsType string `json:"goods_type"`						//商品主类型 :0-虚拟类商品,1-实物类商品
	PassbackParams string `json:"passback_params"`				//公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在同步返回（包括跳转回商户网站）和异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝。
	QuitUrl string `json:"quit_url"`							//用户付款中途退出返回商户网站的地址
	ProductCode string `json:"product_code"`					//销售产品码，商家和支付宝签约的产品码 QUICK_WAP_WAY
	PromoParams string `json:"promo_params"`					//优惠参数 注：仅与支付宝协商后可用
	ExtendParams ExtendParams `json:"extend_params"`			//业务扩展参数
	EnablePayChannels string `json:"enable_pay_channels"`		//可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 注，与disable_pay_channels互斥
	DisablePayChannels string `json:"disable_pay_channels"`		//禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 注，与enable_pay_channels互斥
	StoreId string `json:"store_id"`							//商户门店编号
	SpecifiedChannel string `json:"specified_channel"`			//指定渠道，目前仅支持传入pcredit 若由于用户原因渠道不可用，用户可选择是否用其他渠道支付。 注：该参数不可与花呗分期参数同时传入
	BusinessParams string `json:"business_params"`				//商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	ExtUserInfo ExtUserInfo `json:"ext_user_info"`     			//外部指定买家
}

//业务扩展参数
type ExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id"`	//系统商编号 该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID
	HbFqNum string `json:"hb_fq_num"`								//使用花呗分期要进行的分期数
	HbFqSellerPercent string `json:"hb_fq_seller_percent"`			//使用花呗分期需要卖家承担的手续费比例的百分值，传入100代表100%
	IndustryRefluxInfo string `json:"industry_reflux_info"`			//行业数据回流信息, 详见：地铁支付接口参数补充说明
	CardType string `json:"card_type"`								//卡类型
}

//外部指定买家参数
type ExtUserInfo struct {
	Name string `json:"name"`						//姓名 注： need_check_info=T时该参数才有效
	Mobile string `json:"mobile"`					//手机号 注：该参数暂不校验
	CertType string `json:"cert_type"`				//身份证：IDENTITY_CARD、护照：PASSPORT、军官证：OFFICER_CARD、士兵证：SOLDIER_CARD、户口本：HOKOU等。如有其它类型需要支持，请与蚂蚁金服工作人员联系。 注： need_check_info=T时该参数才有效
	CertNo string `json:"cert_no"`					//证件号 注：need_check_info=T时该参数才有效
	MinAge string `json:"min_age"`					//允许的最小买家年龄，买家年龄必须大于等于所传数值 注：1. need_check_info=T时该参数才有效 2. min_age为整数，必须大于等于0
	FixBuyer string `json:"fix_buyer"`				//是否强制校验付款人身份信息 T:强制校验，F：不强制
	NeedCheckInfo string `json:"need_check_info"`	//是否强制校验身份信息 T:强制校验，F：不强制
}

//下单返回结构体
type AliResTradeWap struct {
	Code string `json:"code"`				//网关返回码,详见文档
	Msg string `json:"msg"`					//网关返回码描述,详见文档
	SubCode string `json:"sub_code"`		//业务返回码，参见具体的API接口文档
	SubMsg string `json:"sub_msg"`			//业务返回码描述，参见具体的API接口文档
	Sign string `json:"sign"`				//签名,详见文档
	OutTradeNo string `json:"out_trade_no"`	//商户网站唯一订单号
	TradeNo string `json:"trade_no"`		//该交易在支付宝系统中的交易流水号。最长64位。
	TotalAmount string `json:"total_amount"`//该笔订单的资金总额，单位为RMB-Yuan。取值范围为[0.01，100000000.00]，精确到小数点后两位。
	SellerId string `json:"seller_id"`		//收款支付宝账号对应的支付宝唯一用户号。 以2088开头的纯16位数字
}

//WAP支付
func (ali AliPay) TradeWap(orderParam AliTradeWapParam) ReturnParam {

	//处理参数
	param,err := ali.tradeWapParam(orderParam)

	if err != nil{
		return ReturnParam{-1,err.Error(),nil}
	}

	//对参数进行排序
	src,_ := tools.SortData(param)

	var hash crypto.Hash
	if orderParam.SignType == "RSA" {
		hash = crypto.SHA1
	} else {
		hash = crypto.SHA256
	}

	privateKey := tools.ParsePrivateKey(ali.config.PrivateKeyString)

	sign, err := ali.SignPKCS1v15([]byte(src), privateKey, hash)

	signWithBase64 := base64.StdEncoding.EncodeToString(sign)

	param.Add("sign", signWithBase64)

	url := AliCommonApi+"?"+param.Encode()

	return ReturnParam{1,"ok",url}
}

//处理交易参数
func (ali AliPay) tradeWapParam(orderParam AliTradeWapParam) (url.Values,error){
	var param = url.Values{}
	param.Add("app_id", ali.config.AppId)
	param.Add("method", "alipay.trade.wap.pay")
	param.Add("format", "JSON")
	param.Add("charset", orderParam.Charset)
	param.Add("sign_type", orderParam.SignType)
	param.Add("timestamp", orderParam.Timestamp)
	param.Add("version", "1.0")

	if orderParam.ReturnUrl != ""{
		param.Add("return_url", orderParam.ReturnUrl)			//
	}
	if orderParam.NotifyUrl != ""{
		param.Add("notify_url", orderParam.NotifyUrl)			//
	}
	if orderParam.AppAuthToken != ""{
		param.Add("app_auth_token", orderParam.AppAuthToken)	//
	}

	extUserInfo := tools.Struct2Map(orderParam.BizContent.ExtUserInfo)

	extendParams := tools.Struct2Map(orderParam.BizContent.ExtendParams)

	bizContent:=tools.Struct2Map(orderParam.BizContent)

	if len(extUserInfo) != 0{
		bizContent["ext_user_info"] = extUserInfo
	}
	if len(extendParams) != 0{
		bizContent["extend_params"] = extendParams
	}

	bizContentString,err := json.Marshal(bizContent)

	if err != nil{
		return nil,errors.New("生成biz_content字符串失败")
	}

	param.Add("biz_content",string(bizContentString))

	return param,nil
}

