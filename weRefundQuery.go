package anypay



//退款订单查询参数
type WeRefundQueryParam struct {
	Appid string `xml:"appid"`					// appid
	MchId string `xml:"mch_id"`					// 商户id
	NonceStr string `xml:"nonce_str"`			// 随机字符串
	Sign string `xml:"sign"`					// 签名
	SignType string `xml:"sign_type"`			// 加密方式
	TransactionId string `xml:"transaction_id"`	// 交易流水号
	OutTradeNo string `xml:"out_trade_no"`		// 商户订单号
	OutRefundNo string `xml:"out_refund_no"`	// 商户退款单号
	RefundId string `xml:"refund_id"`			// 微信退款单号
	Offset string `xml:"offset"`				// 偏移量
}


//退款订单查询返回结构体
type WeResRefundQuery struct {
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


//退款查询
func (w WePay) RefundQuery( queryParam WeRefundQueryParam) ReturnParam {

	return ReturnParam{-1,"not available...",nil}
}

