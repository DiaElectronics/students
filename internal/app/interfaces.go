package app

type Dal interface {
	SaveRate(rate Rate) error
	Rate(CodeFrom, CodeTo string) (Rate, error)
}
