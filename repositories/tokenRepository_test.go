package repositories_test

import (
	"testing"

	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenSuccess(t *testing.T) {
	token, err := repositories.GenerateToken("userid123", "tokenid123", []byte("hmacsecret"), 0)
	assert.Nil(t, err)
	parsedData, _ := repositories.Parse(token, []byte("hmacsecret"))
	assert.Equal(t, float64(0), parsedData["validuntil"])
	assert.Equal(t, "userid123", parsedData["userid"])
}

func TestGenerateTokenWronghmac(t *testing.T) {
	token, err := repositories.GenerateToken("userid123", "tokenid123", []byte("hmacsecret"), 0)
	assert.Nil(t, err)
	_, err = repositories.Parse(token, []byte("wronghmacsecret"))
	assert.NotNil(t, err)
	assert.Equal(t, "signature is invalid", err.Error())

}
