# What is AnyPay？

AnyPay 是一套基于Golang开发的轻量级、易使用的支付组件，目前支持微信支付，支付宝正在开发当中。此组件目前用于我司家博会业务当中，经过实际场景运行稳定可靠。

# How to get it?
go get github.com/panghu1024/anypay

# How to use it?

**微信支付使用参考**

### 1. 生成实例
``` golang
	//初始配置文件
	config := anypay.WeConfig{  
	  AppId:"wx10c46x7s3c78145d", //公众号APPID
	  MchId:"1247551321", //商户ID
	  Key:"zjkk8f123981hj2938123ksd98djfsk12", //商户KEY  
	  CertKeyPath:"./certs/apiclient_key.pem", //证书apiclient_key.pem文件
	  CertP12Path:"./certs/apiclient_cert.p12", //证书apiclient_cert.p12文件
	  CertPemPath:"./certs/apiclient_cert.pem", //证书apiclient_cert.pem文件
	}  

	payment := anypay.NewWePay(config) //创建实例
```

### 2.下单-JSAPI方式支付

``` golang
  
	res := payment.UnifiedOrder(anypay.WeOrderParam{//创建订单
		Body:"测试商品",
		OutTradeNo:"C20190214000133",
		TotalFee:"1",
		SpbillCreateIp:"127.0.0.1", //务必替换成获取的真实IP
		NotifyUrl:"http://notify",
		TradeType:"JSAPI",
		Openid:"obxjktz34If5J6xal0HXGnK23H8E", // JSAPI方式此参数必传
	})

	//结果判断
	if res.Status == 1{//调用成功
		order := res.Data.(anypay.WeResOrder)

		//生成前端支付参数
		resParam := payment.JsApiParam(order.PrepayId)

		if resParam.Status == 1{
			param := resParam.Data.(anypay.WeResJsApi)
			fmt.Println(param)
		}else{
			fmt.Println(res)
		}

	}else{//调用失败
		fmt.Println(res)
	}
```

### 3.下单-NATIVE方式支付

``` golang
  
	res := payment.UnifiedOrder(anypay.WeOrderParam{//创建订单
		Body:"测试商品",
		OutTradeNo:"C20190214000133",
		TotalFee:"1",
		SpbillCreateIp:"127.0.0.1",  //务必替换成获取的真实IP
		NotifyUrl:"http://notify",
		TradeType:"NATIVE",
	})

	//结果判断
	if res.Status == 1{//调用成功
		order := res.Data.(anypay.WeResOrder)
		
		fmt.Println(order.CodeUrl) //微信付款码链接
	}else{//调用失败
		fmt.Println(res)
	}
```

### 4.下单-H5方式支付

``` golang
  
	res := payment.UnifiedOrder(anypay.WeOrderParam{//创建订单
		Body:"测试商品",
		OutTradeNo:"C20190214000133",
		TotalFee:"1",
		SpbillCreateIp:"127.0.0.1",  //务必替换成获取的真实IP
		NotifyUrl:"http://notify",
		TradeType:"MWEB",
	})

	//结果判断
	if res.Status == 1{//调用成功
		order := res.Data.(anypay.WeResOrder)

		fmt.Println(order.MwebUrl) //微信H5支付链接
	}else{//调用失败
		fmt.Println(res)
	}
```

### 5.查询订单

``` golang
  
	res := payment.OrderQuery(anypay.WeOrderQueryParam{//发起查询
		TransactionId:"4200001352401804243426978567",
	})
	
	//结果判断
	if res.Status == 1{//调用成功
		query := res.Data.(WeResOrderQuery)

		fmt.Println(query.TradeState) //打印订单交易状态
	}else{//调用失败
		fmt.Println(res)
	}
```

### 6.订单退款

``` golang
  
	res := payment.Refund(anypay.WeRefundParam{//发起退款
		TransactionId:"4200001352401804243426978567",
		OutRefundNo:"R20190214013007",
		TotalFee:"1",
		RefundFee:"1",
	})

	//结果判断
	if res.Status == 1{//调用成功
		refund := res.Data.(WeResRefund)

		fmt.Println(refund.RefundId) //打印微信退款单号
	}else{//调用失败
		fmt.Println(res)
	}
```

### 7.关闭订单

``` golang
  
	res := payment.CloseOrder(anypay.WeCloseOrderParam{//发起关闭操作
		OutTradeNo:"C20190214000133",
	})

	//结果判断
	if res.Status == 1{//调用成功
		close := res.Data.(WeResCloseOrder)

		fmt.Println(close.ResultCode) //打印状态
	}else{//调用失败
		fmt.Println(res)
	}
```

# 使用反馈以及优化建议
此文档持续更新中,有问题请联系 panghu1024@gmail.com.

# License #

AnyPay is under the terms of the MIT license (see [LICENSE](LICENSE)).
