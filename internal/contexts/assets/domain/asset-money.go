package assetdomain

import (
	"errors"
)

var (
	// ErrMoneyAmountCannotBeNegative represents the error when the money amount is negative.
	ErrMoneyAmountCannotBeNegative = errors.New("money amount cannot be negative")

	// ErrUnsupportedCurrency represents the error when the money currency is invalid.
	ErrUnsupportedCurrency = errors.New("currency not supported, please use USD or EUR")
)

// Money represents the asset money.
type Money struct {
	Amount   float64
	Currency Currency
}

func (m Money) Validate() error {
	if m.Amount < 0 {
		return ErrMoneyAmountCannotBeNegative
	}

	if !m.Currency.IsValid() {
		return ErrUnsupportedCurrency
	}

	return nil
}

// NewMoney creates a new Money object with the given amount and currency.
func NewMoney(amount float64, currency string) (Money, error) {
	money := Money{
		Amount:   amount,
		Currency: Currency(currency),
	}

	err := money.Validate()
	if err != nil {
		return Money{}, err
	}

	return money, nil
}

// Currency returns the money currency code.
type Currency string

const (
	// USD represents the US Dollar currency.
	USD Currency = "USD"

	// EUR represents the Euro currency.
	EUR Currency = "EUR"
)

// String returns the currency code as a string.
func (c Currency) String() string {
	return string(c)
}

// IsValid checks if the currency code is valid.
func (c Currency) IsValid() bool {
	switch c {
	case USD, EUR:
		return true
	}
	return false
}
