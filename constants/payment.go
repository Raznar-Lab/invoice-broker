package constants

type Gateway int
type Currency int

const (
	GATEWAY_XENDIT_ID Gateway = iota + 1
	GATEWAY_PAYPAL_ID
)

const (
	CURRENCY_IDR_ID Currency = iota + 1
	CURRENCY_USD_ID
)

func (g Gateway) String() string {
	if g == GATEWAY_XENDIT_ID {
		return "Xendit"
	}

	if g == GATEWAY_PAYPAL_ID {
		return "Paypal"
	}

	return ""
}

func (c Currency) String() string {
	if c == CURRENCY_IDR_ID {
		return "Rp"
	}

	if c == CURRENCY_USD_ID {
		return "USD$"
	}

	return ""
}


func (g Gateway) StringP(defaultValue string) string {
	if g == GATEWAY_XENDIT_ID {
		return "Xendit"
	}

	if g == GATEWAY_PAYPAL_ID {
		return "Paypal"
	}

	return defaultValue
}

func (c Currency) StringP(defaultValue string) string {
	if c == CURRENCY_IDR_ID {
		return "Rp"
	}

	if c == CURRENCY_USD_ID {
		return "USD$"
	}

	return defaultValue
}
