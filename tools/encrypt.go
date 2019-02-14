package tools

import (
	"encoding/hex"
	"crypto/md5"
	sha12 "crypto/sha1"
)

//Sha1 加密
func Sha1(data string) string{
	sha1 := sha12.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

//MD5加密
func MD5(data string) string{
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
