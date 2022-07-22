package sys_jwt

type (
	UserResult struct {
		Code    int      `json:"code"`
		Data    UserInfo `json:"data"`
		Message string   `json:"msg"`
	}
	UserInfo struct {
		UID   int    `json:"uid"`
		Name  string `json:"name"`
		EName string `json:"ename"`
		Email string `json:"email"`
	}
)
