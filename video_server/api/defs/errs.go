package defs

// 错误信息描述实体
type Err struct{
	ErrorSimple string `json:"error"`  //错误信息提示
	ErrorCode string `json:error_code`  //服务器端自己的error_code
}

// 错误信息返回体
type ErrorResponse struct{
	HttpStatusCode  int  //HTTP状态码
	Error Err  // 错误信息描述实体
}

// 复合结构体需要定义并实例化
var(
	// 请求request参数错误，或request不能被服务端正常解析
	ErroRequestBodyParseFailed = ErrorResponse{HttpStatusCode:400, Error:Err{ErrorSimple:"请求体不能被服务端正常解析.", ErrorCode:"001"}}

	// 用户不存在错误
	ErroUserNotAuthFailed = ErrorResponse{HttpStatusCode:401, Error:Err{ErrorSimple:"请求能被服务端接收，但是用户不存在.", ErrorCode:"002"}}
)


