package constants

type Gateway string
type Currency string

const (
	GATEWAY_XENDIT = Gateway("Xendit")
	GATEWAY_PAYPAL = Gateway("Paypal")
)

const (
	CURRENCY_IDR = Currency("Rp") 
	CURRENCY_USD = Currency("$")
)
