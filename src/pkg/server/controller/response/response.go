package response

// return code
const (
	SUCCESS        = 0
	INVALID_PARAM  = 1000
	INTERNAL_ERROR = 91000
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Render(code int, msg string, data interface{}) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func RenderSuccess(data interface{}) Response {
	return Render(SUCCESS, "success", data)
}

func RenderError(code int, err error) Response {
	return Render(code, err.Error(), nil)
}
