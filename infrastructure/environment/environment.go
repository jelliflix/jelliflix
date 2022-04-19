package environment

import "os"

const (
	key = "ENV"

	DEV  = "DEV"
	PROD = "PROD"
)

func IsDEV() bool {
	return os.Getenv(key) == DEV || os.Getenv(key) == ""
}

func IsProd() bool {
	return os.Getenv(key) == PROD
}
