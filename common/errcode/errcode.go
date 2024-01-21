package errcode

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
	// 鉴权相关
	AUTH_CHECK_ERROR = 10001
	AUTH_EXPIRED     = 10002
	// 用户相关
	USERNAME_EXIST             = 20001
	WRONG_USERNAME_OR_PASSWORD = 20002

	// 业务相关
	NOT_FOUND_SCENE = 30001

)

var msgMap = map[int]string{
	SUCCESS:                    "ok",
	ERROR:                      "fail",
	INVALID_PARAMS:             "请求参数错误",
	AUTH_CHECK_ERROR:           "鉴权失败",
	AUTH_EXPIRED:               "鉴权过期",
	USERNAME_EXIST:             "用户名已存在",
	WRONG_USERNAME_OR_PASSWORD: "用户名或密码错误",
	NOT_FOUND_SCENE:            "未找到场景",
}

// GetMSG 根据code查对应中文信息
func GetMSG(code int) string {
	msg, ok := msgMap[code]
	if !ok {
		return "未知错误"
	}
	return msg
}
