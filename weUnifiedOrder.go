package anypay

import (
	"github.com/panghu1024/anypay/tools"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/xml"
	"strconv"
	"time"
)

//订单参数
type WeOrderParam struct {
	Appid string `xml:"appid"`						// 微信支付分配的公众账号ID（企业号corpid即为此appId）
	MchId string `xml:"mch_id"`						// 微信支付分配的商户号
	DeviceInfo string `xml:"device_info"`			// 自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传"WEB"
	NonceStr string `xml:"nonce_str"`				// 随机字符串，长度要求在32位以内 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
	Sign string `xml:"sign"`						// 通过签名算法计算得出的签名值 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
	SignType string `xml:"sign_type"`				// 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	Body string `xml:"body"`						// 商品简单描述，该字段请按照规范传递 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	Detail string `xml:"detail"`					// 商品详细描述，对于使用单品优惠的商户，改字段必须按照规范上传 - https://pay.weixin.qq.com/wiki/doc/api/danpin.php?chapter=9_102&index=2
	Attach string `xml:"attach"`					// 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用
	OutTradeNo string `xml:"out_trade_no"`			// 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|* 且在同一个商户号下唯一 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	FeeType string `xml:"fee_type"`					// 符合ISO 4217标准的三位字母代码，默认人民币：CNY - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	TotalFee string `xml:"total_fee"`				// 订单总金额，单位为分 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	SpbillCreateIp string `xml:"spbill_create_ip"`	// 支持IPV4和IPV6两种格式的IP地址。调用微信支付API的机器IP
	TimeStart string `xml:"time_start"`				// 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	TimeExpired string `xml:"time_expired"`			// 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。订单失效时间是针对订单号而言的，由于在请求支付的时候有一个必传参数prepay_id只有两小时的有效期，所以在重入时间超过2小时的时候需要重新请求下单接口获取新的prepay_id - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	GoodsTag string `xml:"goods_tag"`				// 订单优惠标记，使用代金券或立减优惠功能时需要的参数 - https://pay.weixin.qq.com/wiki/doc/api/tools/sp_coupon.php?chapter=12_1
	NotifyUrl string `xml:"notify_url"`				// 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数
	TradeType string `xml:"trade_type"`				// JSAPI -JSAPI支付,NATIVE -Native支付,APP -APP支付 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	ProductId string `xml:"product_id"`				// trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义
	LimitPay string `xml:"limit_pay"`				// 上传此参数no_credit--可限制用户不能使用信用卡支付
	Openid string `xml:"openid"`					// trade_type=JSAPI时（即JSAPI支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_4
	Receipt string `xml:"receipt"`					// Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	SceneInfo string `xml:"scene_info"`				// 该字段常用于线下活动时的场景信息上报，支持上报实际门店信息，商户也可以按需求自己上报相关信息。该字段为JSON对象数据，对象格式为{"store_info":{"id": "门店ID","name": "名称","area_code": "编码","address": "地址" }}
}

//下单返回结构体
type WeResOrder struct {
	ReturnCode string `xml:"return_code"`				// SUCCESS/FAIL,此字段是通信标识，表示接口层的请求结果，并非退款状态
	ReturnMsg string `xml:"return_msg"`					// 当return_code为FAIL时返回信息为错误原因
	ResultCode string `xml:"result_code"`				// SUCCESS/FAIL，SUCCESS退款申请接收成功，结果通过退款查询接口查询,FAIL 提交业务失败
	ErrCode string `xml:"err_code"`						// 列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"`				// 结果信息描述
	Appid string `xml:"appid"`							// 微信分配的公众账号ID
	Attach string `xml:"attach"`						// 附加信息
	Body string `xml:"body"`							// 订单BODY
	MchId string `xml:"mch_id"`							// 微信支付分配的商户号
	Detail string `xml:"detail"`						// 订单详细描述信息
	NonceStr string `xml:"nonce_str"`					// 随机字符串，不长于32位
	NotifyUrl string `xml:"notify_url"`					// 异步通知地址
	Openid string `xml:"openid"`						// 用户OpenId
	Sign string `xml:"sign"`							// 签名
	TransactionId string `xml:"transaction_id"`			// 微信订单号
	SpbillCreateIp string `xml:"spbill_create_ip"`		// 用户端IP
	OutTradeNo string `xml:"out_trade_no"`				// 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	TotalFee int `xml:"total_fee"`						// 总金额 单位分
	TradeType string `xml:"trade_type"`					// 交易类型
	MwebUrl string `xml:"mweb_url"`						// H5支付链接 H5支付
	PrepayId string `xml:"prepay_id"`					// 预付码 JSAPI支付
	CodeUrl string `xml:"code_url"`						// 付款码地址 Native支付
}

// JSAPI参数结构体
type WeResJsApi struct {
	TimeStamp string		// 时间戳
	NonceStr string			// 随机字符串
	Package string			// PrepayId 拼接的字符串
	Sign string				// 加密签名
}

// APPAPI参数结构体
type WeResAppApi struct {
	TimeStamp string		// 时间戳
	NonceStr string			// 随机字符串
	Package string			// 固定字符
	PrepayId string			// PrepayId 拼接的字符串
	Sign string				// 加密签名
	PartnerId string		// 商户ID
}

//订单参数检查
func(WePay) paramCheckForUnifiedOrder(orderParam WeOrderParam) bool {
	if orderParam.Body == ""{
		return false
	}
	if orderParam.OutTradeNo == ""{
		return false
	}
	if orderParam.TotalFee == ""{
		return false
	}
	if orderParam.SpbillCreateIp == ""{
		return false
	}
	if orderParam.NotifyUrl == ""{
		return false
	}
	if orderParam.TradeType == ""{
		return false
	}

	return true
}

//统一下单
func (w WePay) UnifiedOrder(orderParam WeOrderParam) ReturnParam {

	if !w.paramCheckForUnifiedOrder(orderParam) {
		return ReturnParam{-1,"请检查订单必传参数",nil}
	}

	nonceStr := tools.GenerateNonceString()

	if orderParam.Appid == ""{
		orderParam.Appid = w.config.AppId	//设置APPID
	}
	orderParam.NonceStr = nonceStr		//设置随机字符串
	orderParam.MchId = w.config.MchId	//设置商户ID

	//订单结构体转MAP
	orderMap := tools.Struct2Map(orderParam)

	//签名数据
	sign,err := tools.GenerateSignString(orderMap,w.config.Key)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	orderParam.Sign = sign//设置签名

	//生成XML
	requestXml,err := tools.GenerateRequestXml(orderParam)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	//发起请求
	r,err := http.Post(ApiUrlMap["UnifiedOrder"],"text/xml",strings.NewReader(requestXml))

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return ReturnParam{-1,err.Error(),nil}
	}

	var res WeResOrder

	xml.Unmarshal([]byte(string(body)),&res)

	if res.ReturnCode != "SUCCESS" || res.ResultCode != "SUCCESS"{
		return ReturnParam{-1,res.ReturnMsg,res}
	}

	return ReturnParam{1,"success",res}
}

//JSAPI 支付参数
func (w WePay) JsApiParam(prepayId string) ReturnParam {
	timeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:10]

	nonceStr := tools.GenerateNonceString()
	packageStr := "prepay_id="+prepayId
	paySign := strings.ToUpper(tools.MD5("appId="+w.config.AppId+"&nonceStr="+nonceStr+"&package="+packageStr+"&signType=MD5&timeStamp="+timeStamp+"&key="+w.config.Key))

	var res WeResJsApi

	res.NonceStr = nonceStr
	res.Package = packageStr
	res.Sign = paySign
	res.TimeStamp = timeStamp

	return ReturnParam{1,"ok",res}
}

//App 支付参数
func (w WePay) AppApiParam(prepayId string) ReturnParam {
	timeStamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:10]

	nonceStr := tools.GenerateNonceString()
	packageStr := "Sign=WXPay"
	paySign := strings.ToUpper(tools.MD5("appid="+w.config.AppId+"&noncestr="+nonceStr+"&package="+packageStr+"&partnerid="+w.config.MchId+"&prepayid="+prepayId+"&timestamp="+timeStamp+"&key="+w.config.Key))

	var res WeResAppApi

	res.NonceStr = nonceStr
	res.Package = packageStr
	res.PrepayId = prepayId
	res.Sign = paySign
	res.PartnerId = w.config.MchId
	res.TimeStamp = timeStamp

	return ReturnParam{1,"ok",res}
}