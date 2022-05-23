package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {

	const min int64 = 0
	const max int64 = 10000

	const n int = 100 // Number of iterations

	for i := 0; i < n; i++ {
		rand := RandomUtils.RandomInt(min, max)
		require.GreaterOrEqual(t, rand, min)
		require.Less(t, rand, max)
	}
}

func TestRandomString(t *testing.T) {

	length := 10

	const n int = 100 // Number of iterations

	for i := 0; i < n; i++ {
		rand := RandomUtils.RandomString(length)
		require.Equal(t, length, len(rand))
	}
}

func TestRandomOwner(t *testing.T) {

	ownerLength := 6

	rand1 := RandomUtils.RandomOwner()

	require.Equal(t, ownerLength, len(rand1))
}

func TestRandomMoney(t *testing.T) {

	const min int64 = 0
	const max int64 = 10000

	const n int = 100 // Number of iterations

	for i := 0; i < n; i++ {
		rand := RandomUtils.RandomMoney()
		require.GreaterOrEqual(t, rand, min)
		require.Less(t, rand, max)
	}
}

func TestRandomCurrency(t *testing.T) {

	currencies := []string{"EUR", "HUF", "USD", "CAD"}

	const n int = 100 // Number of iterations

	for i := 0; i < n; i++ {
		rand := RandomUtils.RandomCurrency()
		require.True(t, isStringPresentInList(rand, currencies))
	}
}

func isStringPresentInList(s string, l []string) bool {

	result := false

	for _, e := range l {
		if e == s {
			result = true
			break
		}
	}

	return result
}
