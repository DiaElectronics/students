package app

import (
	"errors"
	"regexp"
)

// Errors ...
var (
	ErrNegativeAmount    = errors.New("negative money amount")
	ErrCodeTooLong       = errors.New("currency code is too long")
	ErrCodeTooShort      = errors.New("currency code is too short")
	ErrCodeIsEmpty       = errors.New("currency code is empty")
	ErrCodeFormat        = errors.New("currency code format")
	ErrRateIsCloseToZero = errors.New("currency rate can't be so close to zero")
	ErrRateCodesAreSame  = errors.New("rate codes must be different")
)

const (
	currencyCodeRegex  = "^[a-zA-Z]{1}[a-zA-Z\\d]{2}$"
	minimumAllowedRate = 0.00000001
)

type Application struct {
	currencyRegex *regexp.Regexp
	dal           Dal
}

// NewApplication is a constructor
func NewApplication(dal Dal) (*Application, error) {
	compiledRegexp, err := regexp.Compile(currencyCodeRegex)
	if err != nil {
		return nil, err
	}
	return &Application{
		dal:           dal,
		currencyRegex: compiledRegexp,
	}, nil
}

// Exchange returns the amount of money in 'CodeTo' currency
// in the real apps never use float64 for money
func (a *Application) Exchange(CodeFrom, CodeTo string, amount float64) (float64, error) {
	if amount < 0 {
		return 0, ErrNegativeAmount
	}
	if amount == 0 {
		return 0, nil
	}
	errCodeFormat := a.ValidateCurrencyCode([]byte(CodeFrom))
	if errCodeFormat != nil {
		return 0, errCodeFormat
	}
	errCodeFormat = a.ValidateCurrencyCode([]byte(CodeTo))
	if errCodeFormat != nil {
		return 0, errCodeFormat
	}
	if CodeFrom == CodeTo {
		return amount, nil
	}
	needReverse := false
	if CodeFrom > CodeTo {
		needReverse = true
		CodeTo, CodeFrom = CodeFrom, CodeTo
	}

	curRate, err := a.dal.Rate(CodeFrom, CodeTo)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	if !needReverse {
		return amount * curRate.Value, nil
	}

	return amount / curRate.Value, nil
}

// SetRate just sets a new rate for currency pair
func (a *Application) SetRate(rate Rate) error {
	if rate.Value < minimumAllowedRate {
		return ErrRateIsCloseToZero
	}
	if rate.CodeFrom == rate.CodeTo {
		return ErrRateCodesAreSame

	}
	errCodeFormat := a.ValidateCurrencyCode([]byte(rate.CodeFrom))
	if errCodeFormat != nil {
		return errCodeFormat
	}
	errCodeFormat = a.ValidateCurrencyCode([]byte(rate.CodeTo))
	if errCodeFormat != nil {
		return errCodeFormat
	}
	if rate.CodeFrom > rate.CodeTo {
		rate.CodeFrom, rate.CodeTo = rate.CodeTo, rate.CodeFrom
		rate.Value = 1 / rate.Value
	}

	err := a.dal.SaveRate(rate)
	if err != nil {
		return err
	}
	return nil
}

func (a *Application) ValidateCurrencyCode(code []byte) error {
	if len(code) == 0 {
		return ErrCodeIsEmpty
	}
	if len(code) < 3 {
		return ErrCodeTooShort
	}
	if len(code) > 3 {
		return ErrCodeTooLong
	}
	if !a.currencyRegex.Match(code) {
		return ErrCodeFormat
	}
	return nil
}
