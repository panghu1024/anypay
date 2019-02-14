package anypay

//通用返回结构体
type ReturnParam struct {
	Status 	int 			//	状态:1=成功,-1=失败
	Message string			//	提示信息
	Data 	interface{}	//	返回数据
}