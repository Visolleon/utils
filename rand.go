package utils

import (
	"math/rand"
	"time"
)

// Randn 时间戳最为种子的随机数
func Randn(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}

// NewRand 生成随机种子
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
