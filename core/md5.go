package core

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(str string) string {
	md5_ := md5.New()
	md5_.Write([]byte(str))
	cipher_str := md5_.Sum(nil)
	md5_str := hex.EncodeToString(cipher_str)

	return md5_str
}
