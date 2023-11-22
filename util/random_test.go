package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomInt(t *testing.T) {
	randomInt := RandomInt(5, 10)

	require.NotEmpty(t, randomInt)
	require.GreaterOrEqual(t, randomInt, int64(5))
	require.LessOrEqual(t, randomInt, int64(10))
}

func TestRandomString(t *testing.T) {
	randomString := RandomString(10)

	require.NotEmpty(t, randomString)
	require.Len(t, randomString, 10)
}

func TestRandomOwner(t *testing.T) {
	randomOwner := RandomOwner()

	require.NotEmpty(t, randomOwner)
	require.Len(t, randomOwner, 6)
}

func TestRandomMoney(t *testing.T) {
	randomMoney := RandomMoney()

	require.NotEmpty(t, randomMoney)
	require.GreaterOrEqual(t, randomMoney, int64(0))
	require.LessOrEqual(t, randomMoney, int64(1000))
}

func TestRandomCurrency(t *testing.T) {
	randomCurrency := RandomCurrency()

	require.NotEmpty(t, randomCurrency)
	require.Contains(t, supportedCurrencies, randomCurrency)
}
