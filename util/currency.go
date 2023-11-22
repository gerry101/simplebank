package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var supportedCurrencies = []string{USD, EUR, CAD}

func IsSupportedCurrency(currency string) bool {
	for _, supportedCurrency := range supportedCurrencies {
		if currency == supportedCurrency {
			return true
		}
	}

	return false
}
