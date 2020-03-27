package status_error

import (
	"net/http"
)

const (
	// 内部未明确定义的错误
	UnknownError StatusErrorCode = http.StatusInternalServerError*1e6 + 1 + iota
	// 内部用于接收参数时非法的结构
	InvalidStruct
	// Read 调用时发生错误
	ReadFailed
	// 请求参数错误
	InvalidRequestParams
	// HTTP 请求下游失败
	HttpRequestFailed
	// 内部参数错误
	InternalParams
	// 请求 url 错误
	ProtocolParseFailed
	// 请求超时
	RequestTimeout
)

const (
	// http method 不允许
	StatusMethodNotAllowed StatusErrorCode = http.StatusMethodNotAllowed*1e6 + 1 + iota
)

const (
	// @errTalk 请求 JSON 格式错误
	InvalidBodyStruct StatusErrorCode = http.StatusBadRequest*1e6 + 1 + iota
	// @errTalk 请求参数格式错误
	InvalidNonBodyParameters
	// @errTalk 请求参数格式不匹配
	InvalidField
	// @errTalk 解析表单数据失败
	ParseFormFailed
	// @errTalk 读取表单文件数据失败
	ReadFormFileFailed
	// @errTalk 非法签名
	InvalidSecret
)

const (
	// @errTalk 验签失败
	SignFailed StatusErrorCode = http.StatusUnauthorized*1e6 + 1 + iota
)
