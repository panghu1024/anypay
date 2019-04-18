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
		query := res.Data.(anypay.WeResOrderQuery)

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
		refund := res.Data.(anypay.WeResRefund)

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
		close := res.Data.(anypay.WeResCloseOrder)

		fmt.Println(close.ResultCode) //打印状态
	}else{//调用失败
		fmt.Println(res)
	}
```

### 8.企业付款-转账到余额

``` golang
  
	res := payment.TransferBalance(anypay.WeTransferBalanceParam{//发起转账
		PartnerTradeNo:"X20190214100423", // 转账单号 需唯一
		SpbillCreateIp:"127.0.0.1", // 务必替换成用户真实IP
		CheckName:"NO_CHECK",
		Amount:"100",
		Desc:"测试", // 最低转账金额 1CNY
		Openid:"oxst31OVq9B5yB3hYVZxNeKcCG5g",
	})

	if res.Status == 1{//调用成功
		param := res.Data.(anypay.WeResTransferBalance)

		fmt.Println(param)
	}else{//调用失败
		fmt.Println(res)
	}
```

### 9.企业付款-转账到银行卡

``` golang
  
	res := payment.TransferBank(anypay.WeTransferBankParam{//发起转账
		PartnerTradeNo:"X20190214120425", //商户转账单号 需唯一
		EncTrueName:"胖虎", // 账号姓名
		EncBankNo:"622288882878555", //银行卡号
		BankCode:"1001", //银行代码
		Amount:"100", //金额 单位分
		Desc:"测试",
	})

	if res.Status == 1{//调用成功
		param := res.Data.(anypay.WeResTransferBank)

		fmt.Println(param)
	}else{//调用失败
		fmt.Println(res)
	}
```

### 10.接收微信支付异步通知

``` golang
  
	res := payment.Notify(request,true)//接受数据,参数1=http.Request对象,参数2=是否强制验证appId(若存在多个公众号下单统一回调，则不要强制验证)

	if res.Status == 1{//成功获取数据并验签通过
		notifyRes := res.Data.(anypay.WeResNotify)
		fmt.Println(notifyRes)
	}else{
		fmt.Println(res)
	}
```

**支付宝使用参考**

### 1. 生成实例
``` golang
	//初始配置文件
	config := anypay.AliConfig{
		AppId:"your alipay appid",
		PrivateKeyString:"your pem file string",
		PublicKeyString:"your pem file string",
	}

	payment := anypay.NewAliPay(config)//创建实例
```
### 2. 手机网站支付接口2.0
``` golang
	//发起支付
	res := payment.TradeWap(anypay.AliTradeWapParam{
		Charset:"utf-8",//字符集
		SignType:"RSA2",//加密方式，支持RSA,RSA2
		Timestamp:time.Now().Format("2006-01-02 15:04:05"),//时间戳
		NotifyUrl:"domain/Notify",//回调地址
		ReturnUrl:"domain",//返回地址
		BizContent:anypay.TradeWapBizContent{
			Subject:"订单主题",
			OutTradeNo:"O20190313000001",//订单号
			TotalAmount:"0.01",//订单金额
			ProductCode:"QUICK_WAP_WAY",//支付宝产品编码,不可修改
		},
	})

	if res.Status == 1{
		url := res.Data //返回的H5支付链接
	}else{
		fmt.Println(res.Message)
	}
```

### 3. 退款接口
``` golang
	//发起退款
	res := payment.TradeRefund(anypay.AliTradeRefundParam{
		SignType:"RSA2",
		Timestamp:time.Now().Format("2006-01-02 15:04:05"),
		BizContent:anypay.TradeRefundBizContent{
			TradeNo:"2018122522001464891116897117",//交易流水号与订单号二选一
			//OutTradeNo:"T201903010000001",//商户订单号
			RefundAmount:"0.01",
		},
	})
	
	//业务处理
	if res.Status == 1{
		param := res.Data.(anypay.AliResTradeRefund)
		fmt.Println(param)
	}else{
		fmt.Println(res.Message)
	}
```

### 4. 支付通知处理
``` golang
	//发起退款
	res := payment.Notify(request)//接受数据,参数1=http.Request对象
	
	//业务处理
	if res.Status == 1{
		resNotify := res.Data.(anypay.AliResNotify)
		fmt.Println(resNotify)
	}else{
		fmt.Println(res.Message)
	}
```

# Feedback & Suggestion
此文档持续更新中,有问题请联系 panghu1024@gmail.com

# License #

AnyPay is under the terms of the MIT license (see [LICENSE](LICENSE)).
