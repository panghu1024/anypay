package tools

import (
	"math/rand"
	"sort"
	"strings"
	"encoding/xml"
	"reflect"
)

//生成随机字符
func GenerateNonceString() string{
	dict := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var randStr string
	for i:=0;i<32;i++{
		index := rand.Intn(35)
		randStr += dict[index:index+1]
	}

	return randStr
}

//生成验证字符
func GenerateSignString(data map[string]interface{},key string) (str string,err error){

	delete(data,"sign")

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramsStr string
	for _, k := range keys {
		if k == "key"{
			continue
		}
		paramsStr += k+"="+data[k].(string)+"&"
	}

	paramsStr = paramsStr + "key=" + key

	paramsStr = MD5(paramsStr)

	paramsStr = strings.ToUpper(paramsStr)

	return paramsStr,nil
}

//生成XML
func GenerateRequestXml( params interface{} ) (str string,err error){

	data,err := xml.MarshalIndent(&params,""," ")

	if err != nil{
		return "",nil
	}

	return string(data),nil
}

//结构体转MAP
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
    v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {

		if reflect.TypeOf(v.Field(i)).String()=="string" &&  v.Field(i).Interface().(string) == ""{
			continue
		}

		key := SnakeString(t.Field(i).Name)
	 	data[key] = v.Field(i).Interface()
	}

 	return data
}

//驼峰转下划线
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}