package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

// MD5 encodes md5 hash by bytes data
func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// MD5String encodes md5 hash by string data
func MD5String(data string) string {
	return MD5([]byte(data))
}

// MD5Mid encodes md5 hash 16 chars
func MD5Mid(data []byte) string {
	str := MD5(data)
	if len(str) < 30 {
		return ""
	}
	return str[8:24]
}

// MD5MidString encodes md5 hash 16 chars
func MD5MidString(data string) string {
	str := MD5String(data)
	if len(str) < 30 {
		return ""
	}
	return str[8:24]
}

// MD5Data encode ms5 hash for value's json bytes
func MD5Data(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return MD5(b)
}
