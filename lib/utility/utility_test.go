package utility_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestHashPassword(t *testing.T) {
	plain := "strongPassword12"
	hash := utility.HashPassword(plain)
	assert.True(t, utility.ComparePasswords(hash, plain))
}
