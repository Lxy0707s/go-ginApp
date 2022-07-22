package sys_jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-ginApp/src/main/pkg/utils/funcs/hash"
	"sync"
	"time"
)

const (
	SALT = "243223ffslsfsldfl412fdsfsdf" //私钥
)

var (
	once     sync.Once
	instance *jwtTool
)

type (
	jwtTool   struct{}
	SysClaims struct {
		jwt.StandardClaims
		*User
	}
	// User 用于jwt用户认证的功能结构，请不要随意更改
	User struct {
		Uname      string
		Email      string
		Password   string
		Department string
	}

	Secretary struct {
		uName string
		salt  string
	}
	JwtImpl interface {
		GenerateToken(u User) (string, error)
		ParseJwtToken(token string, u string, e ...ErrCheck) (*User, error)
	}
	ErrCheck func(params SysClaims) (bool, error)
)

func NewJwtInstance() *jwtTool {
	if instance == nil {
		once.Do(func() {
			instance = &jwtTool{}
		})
	}
	return instance
}

// GenerateToken 根据传入的Claims和密钥生成规则生成token
func (j *jwtTool) GenerateToken(u User) (string, error) {
	var secretKey, err = j.generateSecretary(u.Uname)
	if err != nil {
		return "", err
	}
	// 设置过期时间
	expireTime := time.Now().Add(time.Hour * 24 * 30)
	//设置jwt当中类似用户对象StandardClaims
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(), //过期时间,int64类型
		IssuedAt:  time.Now().Unix(), //发现人时间,int64类型
		Id:        u.Uname,           //用户UID,字符串类型
		Issuer:    u.Department,      //发行人,字符串类型
	}
	sysClaims := SysClaims{
		StandardClaims: stdClaims,
		User:           &u,
	}
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sysClaims)
	// 使用密钥签名并以字符串形式获取完整的编码令牌
	tokenString, err := token.SignedString([]byte(fmt.Sprint(secretKey)))
	if err != nil {
		return "", err
	}
	// 返回状态
	return tokenString, nil
}

// ParseJwtToken 解析token, 校验token解析是否符合预期
func (j *jwtTool) ParseJwtToken(tokenString string, uName string, errCheck ...ErrCheck) (*User, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}
	var secretKey, err = j.generateSecretary(uName)
	if err != nil {
		return nil, err
	}
	//设置一个空的userStdClaims接受返回来token中携带的信息
	claims := SysClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if claims.User.Uname != uName {
		return nil, errors.New("token parse or user info check error")
	}
	// 自定义检查
	for _, errJudge := range errCheck {
		if ok, err := errJudge(claims); !ok {
			return nil, err
		}
	}
	return claims.User, nil
}

// GenerateSecretary 初始化密钥
func (j *jwtTool) generateSecretary(uName string) (string, error) {
	var secretKey = ""
	if uName != "" {
		secretParams := &Secretary{
			uName: uName,
			salt:  SALT,
		}
		secretKey = hash.MD5Data(secretParams)
	} else {
		return "", errors.New("uName is empty")
	}
	return secretKey, nil
}
