package toolkit

import (
	"crypto/rand"
	"math/big"
)

// 随机数 [0, max) 须 max > 1
func RandInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

const (
	CharsetLowerLetter = "abcdefghijklmnopqrstuvwxyz"
	CharsetUpperLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetLetter      = CharsetLowerLetter + CharsetUpperLetter
	CharsetNumber      = "0123456789"
	CharsetAlpha       = CharsetLetter + CharsetNumber
)

// 随机字符串，自定义字符集
func RandString(size int, charset string) (string, error) {
	if size <= 0 {
		return "", nil
	}

	if len(charset) == 0 {
		charset = CharsetAlpha
	}

	b := make([]byte, size)
	charsetLen := len(charset)
	for i := range b {
		n, err := RandInt(charsetLen)
		if err != nil {
			return "", err
		}
		b[i] = charset[n]
	}

	return string(b), nil
}
