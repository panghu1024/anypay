package anypay

import (
	"crypto"
	"encoding/pem"
	"errors"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
)

//支付宝配置
type AliConfig struct {
	AppId 	string		//APPID
	PrivateKeyString string  //证书key路径
	PublicKeyString string  //证书key路径
}

var AliCommonApi = "https://openapi.alipay.com/gateway.do"

//支付宝结构体
type AliPay struct {
	config AliConfig
}

//创建支付实例
func NewAliPay(config AliConfig) (AliPay){
	if config.AppId == ""{
		panic("AppId不可为空")
	}

	if config.PrivateKeyString == ""{
		panic("支付宝私钥不可为空")
	}

	if config.PublicKeyString == ""{
		panic("支付宝公钥不可为空")
	}

	return AliPay{config:config}
}

//加密
func (ali AliPay) SignPKCS1v15(src, key []byte, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

//验证
func (ali AliPay) VerifyPKCS1v15(src []byte, sign []byte, key []byte, hash crypto.Hash) (bool,error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return false,errors.New("public key error")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false,err
	}

	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), hash, hashed, sign)

	if err != nil{
		return false,err
	}

	return true,nil
}
