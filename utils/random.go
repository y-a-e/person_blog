package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// GenerateVerificationCode 生成一个指定长度的随机验证码
func GenerateVerificationCode(length int) string {
	//time.Now().UnixNano()表示当前时间的纳秒级时间戳
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//math.Pow10(length)表示10的length次方
	//r.Intn表示生成一个0到math.Pow10(length)-1之间的随机整数
	//"%0*d"，未满"length"位数，在前面自动补0，
	return fmt.Sprintf("%0*d", length, r.Intn(int(math.Pow10(length))))
}
