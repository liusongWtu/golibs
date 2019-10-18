package utils

import (
	"math/rand"
	"time"
)

func init() {
	//以时间作为初始化种子
	rand.Seed(time.Now().UnixNano())
}

// RandInt64 get the random numer in [min, max]
func RandInt64(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	x := rand.Int63n(max-min) + min
	return x
}

// RandInt get the random numer in [min, max]
func RandInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	x := rand.Intn(max-min) + min
	return x
}

//乱序
func Perm(n int) []int {
	return rand.Perm(n)
}
