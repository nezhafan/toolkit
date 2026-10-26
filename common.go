package toolkit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strings"
	"unicode"
)

// 打印JSON
func JSON(v any) string {
	bs, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bs)
}

// 是否为ASCII
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// 是否有效数字（float64存在NaN和Inf的可能）
func IsNumber(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}

// 除法，截断保留x位小数 (非四舍五入)
func DivTrunc[E Number](a, b E, precision uint8) float64 {
	if a == 0 || b == 0 {
		return 0
	}
	return Trunc(float64(a)/float64(b), precision)
}

// 截断保留x位小数 (非四舍五入)
func Trunc(n float64, precision uint8) float64 {
	pow := math.Pow(10, float64(precision))
	if n = math.Trunc(n * pow); !IsNumber(n) {
		return 0
	}
	return n / pow
}

// 加密
func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 加密
func HmacSha256(s string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 追踪调用位置
func Callers() []string {
	pc := make([]uintptr, 8)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return nil
	}

	callers := make([]string, 0, n)

	pc = pc[:n]
	for _, p := range pc {
		fn := runtime.FuncForPC(p)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(p)
		arr := strings.Split(file, "/")
		if len(arr) > 3 {
			file = strings.Join(arr[len(arr)-3:], "/")
		}
		callers = append(callers, fmt.Sprintf("%s:%d", file, line))
	}

	return callers
}
