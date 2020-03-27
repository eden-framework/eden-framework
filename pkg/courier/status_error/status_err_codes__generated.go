package status_error

func init() {
	StatusErrorCodes.Register("InvalidBodyStruct", 400000001, "请求 JSON 格式错误", "", true)
	StatusErrorCodes.Register("InvalidNonBodyParameters", 400000002, "请求参数格式错误", "", true)
	StatusErrorCodes.Register("InvalidField", 400000003, "请求参数格式不匹配", "", true)
	StatusErrorCodes.Register("ParseFormFailed", 400000004, "解析表单数据失败", "", true)
	StatusErrorCodes.Register("ReadFormFileFailed", 400000005, "读取表单文件数据失败", "", true)
	StatusErrorCodes.Register("InvalidSecret", 400000006, "非法签名", "", true)
	StatusErrorCodes.Register("SignFailed", 401000001, "验签失败", "", true)
	StatusErrorCodes.Register("StatusMethodNotAllowed", 405000001, "http method 不允许", "", false)
	StatusErrorCodes.Register("UnknownError", 500000001, "内部未明确定义的错误", "", false)
	StatusErrorCodes.Register("InvalidStruct", 500000002, "内部用于接收参数时非法的结构", "", false)
	StatusErrorCodes.Register("ReadFailed", 500000003, "Read 调用时发生错误", "", false)
	StatusErrorCodes.Register("InvalidRequestParams", 500000004, "请求参数错误", "", false)
	StatusErrorCodes.Register("HttpRequestFailed", 500000005, "HTTP 请求下游失败", "", false)
	StatusErrorCodes.Register("InternalParams", 500000006, "内部参数错误", "", false)
	StatusErrorCodes.Register("ProtocolParseFailed", 500000007, "请求 url 错误", "", false)
	StatusErrorCodes.Register("RequestTimeout", 500000008, "请求超时", "", false)
}
