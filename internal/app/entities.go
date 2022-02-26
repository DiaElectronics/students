package app

type Currency struct {
	Code string
}

type Rate struct {
	CodeFrom string
	CodeTo   string
	Value    float64
}
