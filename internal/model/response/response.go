package response

// 错误码
const (
	CodeSuccess = 0
	//TODO:
)

// Response 请求回复通用结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty" `
}

// Success ...
func Success(msg string, data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: msg,
		Data:    data,
	}
}

// Fail ...
func Fail(code int, err error) *Response {
	return &Response{
		Code:    code,
		Message: err.Error(),
	}
}
