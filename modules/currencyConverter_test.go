package modules_test

import (
	"testing"

	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/modules"
	"github.com/stretchr/testify/assert"
)

func TestConvertCurrency(t *testing.T) {
	ccm := modules.NewTestCCM()
	mu1 := models.MonetaryUnit{Amount: 123, Currency: "AED"}
	mu1Converted := ccm.ConvertCurrency(mu1, "EUR")
	assert.Equal(t, models.MonetaryUnit{Amount: 29.1018, Currency: "EUR"}, mu1Converted)
}

func TestConvertSameCurrency(t *testing.T) {
	ccm := modules.NewTestCCM()
	mu1 := models.MonetaryUnit{Amount: 123, Currency: "AED"}
	mu1Converted := ccm.ConvertCurrency(mu1, "AED")
	assert.Equal(t, models.MonetaryUnit{Amount: 123, Currency: "AED"}, mu1Converted)
}
