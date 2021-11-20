package utility_test

import (
	"testing"

	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	plain := "strongPassword12"
	hash := utility.HashPassword(plain)
	assert.True(t, utility.ComparePasswords(hash, plain))
}
