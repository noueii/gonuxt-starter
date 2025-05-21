package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomFloat(min, max float64) float64 {
	return (max-min)*rand.Float64() + min
}

func RandomString(n int) string {
	var sb strings.Builder

	alpha := "abcdefghijklmnopqrstuvwxyz"

	for range n {
		sb.WriteByte(alpha[rand.Intn(len(alpha))])
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}
