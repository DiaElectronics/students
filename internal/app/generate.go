package app

//go:generate mockgen -source=interfaces.go -destination=testing.generated.go -package app -self_package=usd_converter/internal/app Dal
