package errors

import (
	"log/slog"
)

func PanicOnErr(err error) {
	if err != nil {
		slog.Error("Error occurred: ", err)
		panic(err)
	}
}

func FatalOnErr(err error) {
	if err != nil {
		slog.Error("Error: %s", err)
	}
}

// the current idea with this function is that not every error i encounter should end the program
func LogIfError(err error) bool {
	if err != nil {
		slog.Error("Error occurred: ", err)
		return true
	}
	return false
}
