package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-ginApp/src/main/pkg/common/middleware/sys_jwt"
	"time"
)

var AppIss = "apm-center"

type (
	UserOp struct {
		jwt sys_jwt.JwtImpl
	}
)

func NewInstance() *UserOp {
	return &UserOp{
		jwt: sys_jwt.NewJwtInstance(),
	}
}

// Demo 用户注册
func (u *UserOp) Demo(s *sys_jwt.User) {
	//用下面自定义 claim
	var uerClaims = sys_jwt.User{
		Uname:      "admin",
		Password:   "123456",
		Department: AppIss,
	}
	// 生成token
	token, err := u.jwt.GenerateToken(uerClaims)
	if err != nil {
		return
	}
	fmt.Println("generate token:", token)
	// 解析token
	jwtToken, err := u.jwt.ParseJwtToken(token, "admin.com", u.jwtCheck)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("token parse:", jwtToken)
}

// 注册生成token
func (u *UserOp) Register(sUser *sys_jwt.User) (string, error) {
	if sUser == nil {
		return "", errors.New("user info is nil")
	}
	// 生成token
	sUser.Department = AppIss
	token, err := u.jwt.GenerateToken(*sUser)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Login 登录检查
func (u *UserOp) Login(s *sys_jwt.User) (string, error) {
	// 查询用户是否存在，返回用户基本信息

	// 校验用户token是否匹配

	// 返回用户token进行后续处理 / 返回错误

	return "", nil
}

// 需要校验的字段，返回校验结果和报错信息
func (u *UserOp) jwtCheck(params sys_jwt.SysClaims) (bool, error) {
	var jwtParams, err = json.Marshal(params)
	var jwtClaims *sys_jwt.SysClaims
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(jwtParams, &jwtClaims)
	if err != nil || jwtClaims == nil {
		return false, err
	}
	//检查一下超时没有,超时就返回错误
	if jwtClaims.VerifyExpiresAt(int64(time.Now().Unix()), true) == false {
		return false, errors.New("超时")
	}
	//检查一下发行方正确,令牌的发行方错误
	if !jwtClaims.VerifyIssuer(AppIss, true) {
		return false, errors.New("token's issuer is wrong")
	}
	return true, nil
}
