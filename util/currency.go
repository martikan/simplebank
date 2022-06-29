package util

// Constants for all supported currencies.
const (
	USD = "USD"
	EUR = "EUR"
	HUF = "HUF"
	CAD = "CAD"
)

var (
	CurrencyUtils currencyUtilsInterface = &currencyUtils{}
)

type currencyUtils struct{}

type currencyUtilsInterface interface {
	IsSupportedCurrency(string) bool
}

// IsSupportedCurrency returns true if the specified currency is supported.
func (c *currencyUtils) IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, HUF, CAD:
		return true
	default:
		return false
	}
}
