package sys_jwt

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go-ginApp/src/main/pkg/utils/base_struct"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// 获取服务器用户信息
func GetUserRemoteInfo(ticket string, sConfig base_struct.JwtConfig) (info *UserInfo) {
	ssoLoginURI := "/account/user"
	staffURL := sConfig.JwtUrl + ssoLoginURI
	params := map[string]string{
		"ticket": ticket,
		"appId":  sConfig.AppId,
	}
	params["sign"] = Sign(params, sConfig)
	body, errJ := json.Marshal(params)
	if errJ != nil {
		log.Println("json error!" + errJ.Error())
		return nil
	}
	//http send
	req, errN := http.NewRequest("POST", staffURL, bytes.NewReader(body))
	if errN != nil {
		log.Println("SSOGetInfo unmarshal error!")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		//Transport: transport,
		Timeout: 6 * time.Second,
	}
	resp, errD := client.Do(req)
	if errD != nil || resp == nil {
		log.Println("client do error!" + "err msg: " + errD.Error())
		return nil
	}
	defer resp.Body.Close()
	//io.Reader
	body, errR := ioutil.ReadAll(resp.Body)
	if errR != nil {
		log.Println("io read error!" + errR.Error())
	}
	var result UserResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Println("SSOGetInfo unmarshal error!")
		return nil
	}
	if result.Code != 0 {
		log.Println("code != 0!", "code", result.Code, "msg", result.Message)
		return nil
	}
	return &result.Data
}

// 获取jwtToken
func GetJwtToken(sConfig base_struct.JwtConfig, params map[string]interface{}) string {
	claims := make(jwt.MapClaims)
	claims["name"] = params["name"]
	claims["uid"] = params["uid"]
	claims["email"] = params["email"]
	claims["nbf"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(sConfig.AppSecret))

	if err != nil {
		log.Println("get token error :", err.Error())
		return ""
	}
	return tokenString
}

// 解析token
func ParseToken(token string, config base_struct.JwtConfig) (checkMsg string, claims map[string]interface{}) {
	claims = make(map[string]interface{})
	if token == "" {
		return "token 不可为空", claims
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppSecret), nil
	})

	if err != nil || !jwtToken.Valid {
		return "token 解析错误", claims
	}

	var jOk bool
	claims, jOk = jwtToken.Claims.(jwt.MapClaims)
	if !jOk {
		return "token 二段解析错误", claims
	}

	_, ok := claims["name"].(string)
	if !ok || claims["name"] == "" {
		return "解析后无法获取用户名", claims
	}
	if !ok || claims["uid"] == "" {
		return "解析后无法获取用户ID", claims
	}

	return "", claims
}

// 签名
func Sign(params map[string]string, sConfig base_struct.JwtConfig) string {

	signStr := sliceOfKeys(params)
	h := sha1.New()
	h.Write([]byte(signStr))
	hashBytes := h.Sum(nil)
	sha1Str := hex.EncodeToString(hashBytes) + string(sConfig.AppSecret) + sConfig.AppSalt
	md := md5.New()
	md.Write([]byte(sha1Str))
	result := hex.EncodeToString(md.Sum(nil))
	return result[5:29]
}

func sliceOfKeys(params map[string]string) string {
	var keys []string
	var pairs []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		pairs = append(pairs, k+"="+params[k])
	}
	return strings.Join(pairs, "&")
}
