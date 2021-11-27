package modules

import (
	"encoding/json"
	"fmt"

	"github.com/Kaibling/psychic-octo-stock/lib/config"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
)

var CCM *CurrencyConverterModule

type currencyRate struct {
	Code        string
	AlphaCode   string
	NumericCode string
	Name        string
	Rate        float64
	Date        string
	InverseRate float64
}

type CurrencyConverterModule struct {
	BaseCurrency string
	Rates        map[string]map[string]currencyRate
}

func NewTestCCM() *CurrencyConverterModule {

	rates := map[string]map[string]currencyRate{
		"EUR": {"AED": currencyRate{Code: "AED",
			AlphaCode:   "AED",
			NumericCode: "784",
			Name:        "U.A.E Dirham",
			Rate:        4.2264,
			Date:        "Sat, 27 Nov 2021 11:55:01 GMT",
			InverseRate: 0.2366}},
		"AED": {"EUR": currencyRate{Code: "EUR",
			AlphaCode:   "EUR",
			NumericCode: "785",
			Name:        "EURO",
			Rate:        0.2366,
			Date:        "Sat, 27 Nov 2021 11:55:01 GMT",
			InverseRate: 4.2264}}}

	return &CurrencyConverterModule{Rates: rates, BaseCurrency: "EUR"}

}
func SetGlobalCCM(ccm *CurrencyConverterModule) {
	CCM = ccm
}

func NewCCM() *CurrencyConverterModule {
	loadedRates := initRates()
	return &CurrencyConverterModule{Rates: loadedRates}
}

func (s *CurrencyConverterModule) ConvertCurrency(mu models.MonetaryUnit, currency string) models.MonetaryUnit {
	if mu.Currency == currency {
		//skip same currency
		return mu
	}
	exchangeRate := s.Rates[mu.Currency][currency].Rate
	return models.MonetaryUnit{Amount: mu.Amount * exchangeRate, Currency: currency}
}

func (s *CurrencyConverterModule) SupportedCurrency(currency string) bool {
	for k := range s.Rates {
		if k == currency {
			return true
		}
	}
	return false

}

func (s *CurrencyConverterModule) AddAndConvertFunds(mu1 models.MonetaryUnit, mu2 models.MonetaryUnit) models.MonetaryUnit {
	//if both have the same currency, just add it
	//if one has the default one, take the default
	//if nobody has default, convert both to default
	if mu1.Currency == mu2.Currency {
		return models.MonetaryUnit{Amount: mu1.Amount + mu2.Amount, Currency: mu1.Currency}
	} else if mu1.Currency == config.Config.Currency || mu2.Currency == config.Config.Currency {
		if mu1.Currency == config.Config.Currency {
			mu2Converted := s.ConvertCurrency(mu2, config.Config.Currency)
			return models.MonetaryUnit{Amount: mu1.Amount + mu2Converted.Amount, Currency: config.Config.Currency}
		}
		mu1Converted := s.ConvertCurrency(mu1, config.Config.Currency)
		return models.MonetaryUnit{Amount: mu2.Amount + mu1Converted.Amount, Currency: config.Config.Currency}
	}
	mu1Converted := s.ConvertCurrency(mu1, config.Config.Currency)
	mu2Converted := s.ConvertCurrency(mu2, config.Config.Currency)
	return models.MonetaryUnit{Amount: mu1Converted.Amount + mu2Converted.Amount, Currency: config.Config.Currency}
}

func (s *CurrencyConverterModule) SubtractAndConvertFunds(mu1 models.MonetaryUnit, mu2 models.MonetaryUnit) models.MonetaryUnit {
	//if both have the same currency, just add it
	//if one has the default one, take the default
	//if nobody has default, convert both to default
	if mu1.Currency == mu2.Currency {
		return models.MonetaryUnit{Amount: mu1.Amount - mu2.Amount, Currency: mu1.Currency}
	} else if mu1.Currency == config.Config.Currency || mu2.Currency == config.Config.Currency {
		if mu1.Currency == config.Config.Currency {
			mu2Converted := s.ConvertCurrency(mu2, config.Config.Currency)
			return models.MonetaryUnit{Amount: mu1.Amount - mu2Converted.Amount, Currency: config.Config.Currency}
		}
		mu1Converted := s.ConvertCurrency(mu1, config.Config.Currency)
		return models.MonetaryUnit{Amount: mu2.Amount - mu1Converted.Amount, Currency: config.Config.Currency}
	}
	mu1Converted := s.ConvertCurrency(mu1, config.Config.Currency)
	mu2Converted := s.ConvertCurrency(mu2, config.Config.Currency)
	return models.MonetaryUnit{Amount: mu1Converted.Amount - mu2Converted.Amount, Currency: config.Config.Currency}
}

func initRates() map[string]map[string]currencyRate {
	//Currencies := []string{"USD", "EUR", "GBP", "AUD", "CAD", "JPY", "CHF", "KMF", "AFN", "ALL", "DZD", "AOA", "ARS", "AMD", "AWG", "AZN", "BSD", "BHD", "BDT", "BBD", "BYR", "BYN", "BZD", "BOB", "BAM", "BWP", "BRL", "BND", "BGN", "BIF", "KHR", "CVE", "XAF", "XPF", "CLP", "CNY", "COP", "CDF", "CRC", "HRK", "CUP", "CZK", "DKK", "DJF", "DOP", "XCD", "EGP", "ERN", "ETB", "FJD", "GMD", "GEL", "GHS", "GIP", "GTQ", "GNF", "GYD", "HTG", "HNL", "HKD", "HUF", "ISK", "INR", "IDR", "IRR", "IQD", "ILS", "JMD", "JOD", "KZT", "KES", "KWD", "KGS", "LAK", "LVL", "LBP", "LSL", "LRD", "LYD", "LTL", "MOP", "MKD", "MGA", "MWK", "MYR", "MVR", "MRO", "MRU", "MUR", "MXN", "MDL", "MNT", "MAD", "MZN", "MMK", "NAD", "NPR", "ANG", "TWD", "TMT", "NZD", "NIO", "NGN", "NOK", "OMR", "PKR", "PAB", "PGK", "PYG", "PEN", "PHP", "PLN", "QAR", "RON", "RUB", "RWF", "SVC", "WST", "STN", "SAR", "RSD", "SCR", "SLL", "SGD", "SBD", "SOS", "ZAR", "KRW", "SSP", "LKR", "SDG", "SRD", "SZL", "SEK", "SYP", "TJS", "TZS", "THB", "TOP", "TTD", "TND", "TRY", "AED", "UGX", "UAH", "UYU", "UZS", "VUV", "VES", "VEF", "VND", "XOF", "YER", "ZMW"}
	Currencies := []string{"USD", "EUR"}
	currencyRate := loadFloatRates(Currencies)
	return currencyRate
}

func loadFloatRates(Currencies []string) map[string]map[string]currencyRate {

	returnRates := make(map[string]map[string]currencyRate)
	for _, v := range Currencies {
		url := fmt.Sprintf("http://www.floatrates.com/daily/%s.json", v)
		body := utility.GetRequest(url)
		var rates map[string]currencyRate
		err := json.Unmarshal(body, &rates)
		if err != nil {
			fmt.Println(err)
			continue
		}
		returnRates[v] = rates
	}
	return returnRates

}
