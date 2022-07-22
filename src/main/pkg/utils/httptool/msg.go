package httptool

var MsgFlags = map[int]string{
	OK:            "ok",
	SUCCESS:       "success",
	ERROR:         "fail",
	InvalidParams: "请求异常",
	ErrNotFound:   "对应数据不存在",
	ErrorDb:       "数据库异常",

	ErrorAPIAuthCheckTokenFail:    "API Token鉴权失败",
	ErrorAPIAuthCheckTokenTimeout: "API Token已超时",
	ErrorAPIConfig:                "Current permission limit configuration exception",
	ErrorAPIAuth:                  "The current incoming APP Token is incorrect",
	EmptyAPIAuth:                  "The current parameter is abnormal, and the APP Token parameter is missing",
	ErrorAPILimit:                 "The current API request frequency exceeds the limit, please try again later",

	ErrorSSOValidate: "用户Token错误",

	LogTypeMysqlAdd:    "MysqlAdd",
	LogTypeMysqlUpdate: "MysqlUpdate",
	LogTypeMysqlDelete: "MysqlDelete",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
