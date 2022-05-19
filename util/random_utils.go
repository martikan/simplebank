package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

var (
	RandomUtils randomUtilsInterface = &randomUtils{}
)

type randomUtils struct{}

type randomUtilsInterface interface {
	RandomInt(int64, int64) int64
	RandomString(int) string
	RandomOwner() string
	RandomMoney() int64
	RandomCurrency() string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func (r *randomUtils) RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func (r *randomUtils) RandomString(n int) string {

	k := len(alphabet)

	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func (r *randomUtils) RandomOwner() string {
	return RandomUtils.RandomString(6)
}

// RandomMoney generates a random amount of money
func (r *randomUtils) RandomMoney() int64 {
	return RandomUtils.RandomInt(0, 10000)
}

// RandomCurrency generates a random currency code
func (r *randomUtils) RandomCurrency() string {
	currencies := []string{"EUR", "HUF", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
