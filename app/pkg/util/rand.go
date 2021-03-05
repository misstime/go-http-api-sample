package util

import (
	"math/rand"
	"strings"
	"time"
)

// GenerateRandomDigits 生成指定长度的随机字符串，每个字符都是 0-9 之间的数字。
// 例：GenerateRandomDigits(6) -> 123456
func GenerateRandomDigits(width uint8) string {
	digits := [10]string{`0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`}
	rand.Seed(time.Now().UnixNano())
	var (
		i  uint8
		sb strings.Builder
	)
	for i = 0; i < width; i++ {
		sb.WriteString(digits[rand.Intn(10)])
	}
	return sb.String()
}
