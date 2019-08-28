package util

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func GetMd5(str string) (strMd5 string) {
	data := []byte(str)
	has := md5.Sum(data)
	strMd5 = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return
}

const numberBytes = "0123456789012345678901234567890123456789012345678901234567898"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var srcLogin = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, srcLogin.Int63(), letterIdxMax; i >= 0; {
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandNumber(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, srcLogin.Int63(), letterIdxMax; i >= 0; {
		if idx := int(cache & letterIdxMask); idx < len(numberBytes) {
			b[i] = numberBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
