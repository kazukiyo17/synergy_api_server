package errcode

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_CODE_PARAMS_ERROR = 10001

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
	ERROR_AUTH_USER_EXIST          = 20005
	ERROR_AUTH_EXPIRED             = 20006
)

var msgMap = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_CODE_PARAMS_ERROR:        "请求参数错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
}

// GetMSG 根据code查对应中文信息
func GetMSG(code int) string {
	msg, ok := msgMap[code]
	if !ok {
		return "未知错误"
	}
	return msg
}
