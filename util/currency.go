package util

func allValidCurrency() []string {
	// add validCurrency when needed
	validCurrency := []string{
		"USD", "YEN", "RMB",
	}
	return validCurrency
}

// check if args currency is supported in this system or not
func IsSupportedCurrency(currency string) bool {
	for _, c := range allValidCurrency() {
		if c == currency {
			return true
		}
	}
	return false
}
