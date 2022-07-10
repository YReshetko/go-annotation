package api

type Location struct {
	Country     string
	CountryCode *string
	Region      string
	City        string
	Currencies  map[Currency]ExchangeRate
}

type Currency string

const (
	USD Currency = "USD"
	GBP Currency = "GBP"
	EUR Currency = "EUR"
)

type ExchangeRate struct {
	RoundingMode int
}
