package sys_jwt

import "go-ginApp/src/main/pkg/utils/base_struct"

func Login(token, ticket string, config base_struct.JwtConfig) map[string]interface{} {
	checkMsg, _ := ParseToken(token, config)
	result := map[string]interface{}{
		"code":    0,
		"message": "ok",
	}
	if checkMsg != "" {
		if ticket == "" {
			result["code"] = 1002
			result["message"] = "Failed"
			return result
		}

		info := GetUserRemoteInfo(ticket, config)
		if info == nil {
			result["code"] = 1002
			result["message"] = "Failed"
			return result
		}
		result["name"] = info.Name
		result["uid"] = info.UID
		result["email"] = info.Email
		result["jwtToken"] = GetJwtToken(config, result)
	}
	return result
}
