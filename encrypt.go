package utils

import (
	"crypto/md5"
	"fmt"
)

// MD5Encode MD5加密
func MD5Encode(pwd, salt string) []byte {
	if len(salt) > 0 {
		pwd = fmt.Sprintf("%s-%s", pwd, salt)
	}
	m := md5.New()
	m.Write([]byte(pwd))
	return m.Sum(nil)
}

// MD5EncodeStr  获取 md5 编码密码字符串
func MD5EncodeStr(pwd string) string {
	res := MD5Encode(pwd, "")
	if res != nil {
		return fmt.Sprintf("%x", res)
	}
	return ""
}
