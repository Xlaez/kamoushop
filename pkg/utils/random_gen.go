package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alp = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomStr(n int) string {
	var sb strings.Builder

	k := len(alp)

	for i := 0; i < n; i++ {
		c := alp[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
