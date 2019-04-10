package anypay

import (
	"github.com/panghu1024/anypay/tools"
	"crypto"
	"net/http"
	"encoding/base64"
)


//通知返回结构体
type AliResNotify struct {
	AuthAppId         string `json:"auth_app_id"`         // App Id
	NotifyTime        string `json:"notify_time"`         // 通知时间
	NotifyType        string `json:"notify_type"`         // 通知类型
	NotifyId          string `json:"notify_id"`           // 通知校验ID
	AppId             string `json:"app_id"`              // 开发者的app_id
	Charset           string `json:"charset"`             // 编码格式
	Version           string `json:"version"`             // 接口版本
	SignType          string `json:"sign_type"`           // 签名类型
	Sign              string `json:"sign"`                // 签名
	TradeNo           string `json:"trade_no"`            // 支付宝交易号
	OutTradeNo        string `json:"out_trade_no"`        // 商户订单号
	OutBizNo          string `json:"out_biz_no"`          // 商户业务号
	BuyerId           string `json:"buyer_id"`            // 买家支付宝用户号
	BuyerLogonId      string `json:"buyer_logon_id"`      // 买家支付宝账号
	SellerId          string `json:"seller_id"`           // 卖家支付宝用户号
	SellerEmail       string `json:"seller_email"`        // 卖家支付宝账号
	TradeStatus       string `json:"trade_status"`        // 交易状态
	TotalAmount       string `json:"total_amount"`        // 订单金额
	ReceiptAmount     string `json:"receipt_amount"`      // 实收金额
	InvoiceAmount     string `json:"invoice_amount"`      // 开票金额
	BuyerPayAmount    string `json:"buyer_pay_amount"`    // 付款金额
	PointAmount       string `json:"point_amount"`        // 集分宝金额
	RefundFee         string `json:"refund_fee"`          // 总退款金额
	Subject           string `json:"subject"`             // 总退款金额
	Body              string `json:"body"`                // 商品描述
	GmtCreate         string `json:"gmt_create"`          // 交易创建时间
	GmtPayment        string `json:"gmt_payment"`         // 交易付款时间
	GmtRefund         string `json:"gmt_refund"`          // 交易退款时间
	GmtClose          string `json:"gmt_close"`           // 交易结束时间
	FundBillList      string `json:"fund_bill_list"`      // 支付金额信息
	PassbackParams    string `json:"passback_params"`     // 回传参数
	VoucherDetailList string `json:"voucher_detail_list"` // 优惠券信息
}

//支付通知
func (ali AliPay) Notify(req *http.Request) ReturnParam {

	notify,_ := ali.notifyParam(req)

	publicKey := tools.ParsePublicKey(ali.config.PublicKeyString)

	req.Form.Del("sign")
	req.Form.Del("sign_type")

	//对参数进行排序
	src,_ := tools.SortData(req.Form)

	var hash crypto.Hash
	if notify.SignType == "RSA" {
		hash = crypto.SHA1
	} else {
		hash = crypto.SHA256
	}

	signBytes, err := base64.StdEncoding.DecodeString(notify.Sign)
	if err != nil {
		return ReturnParam{-1,"base64解码签名失败",err.Error()}
	}

	res,err := ali.VerifyPKCS1v15([]byte(src),signBytes,[]byte(publicKey),hash)

	if res == false{
		return ReturnParam{-1,"fail",err.Error()}
	}

	return ReturnParam{1,"ok",notify}
}

//参数整合
func (ali AliPay) notifyParam(req *http.Request) (AliResNotify,error){

	notify := AliResNotify{}
	notify.AppId = req.FormValue("app_id")
	notify.AuthAppId = req.FormValue("auth_app_id")
	notify.NotifyId = req.FormValue("notify_id")
	notify.NotifyType = req.FormValue("notify_type")
	notify.NotifyTime = req.FormValue("notify_time")
	notify.TradeNo = req.FormValue("trade_no")
	notify.TradeStatus = req.FormValue("trade_status")
	notify.TotalAmount = req.FormValue("total_amount")
	notify.ReceiptAmount = req.FormValue("receipt_amount")
	notify.InvoiceAmount = req.FormValue("invoice_amount")
	notify.BuyerPayAmount = req.FormValue("buyer_pay_amount")
	notify.SellerId = req.FormValue("seller_id")
	notify.SellerEmail = req.FormValue("seller_email")
	notify.BuyerId = req.FormValue("buyer_id")
	notify.BuyerLogonId = req.FormValue("buyer_logon_id")
	notify.FundBillList = req.FormValue("fund_bill_list")
	notify.Charset = req.FormValue("charset")
	notify.PointAmount = req.FormValue("point_amount")
	notify.OutTradeNo = req.FormValue("out_trade_no")
	notify.OutBizNo = req.FormValue("out_biz_no")
	notify.GmtCreate = req.FormValue("gmt_create")
	notify.GmtPayment = req.FormValue("gmt_payment")
	notify.GmtRefund = req.FormValue("gmt_refund")
	notify.GmtClose = req.FormValue("gmt_close")
	notify.Subject = req.FormValue("subject")
	notify.Body = req.FormValue("body")
	notify.RefundFee = req.FormValue("refund_fee")
	notify.Version = req.FormValue("version")
	notify.SignType = req.FormValue("sign_type")
	notify.Sign = req.FormValue("sign")
	notify.PassbackParams = req.FormValue("passback_params")
	notify.VoucherDetailList = req.FormValue("voucher_detail_list")

	return notify,nil
}
