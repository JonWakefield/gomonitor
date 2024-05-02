package errors

import (
	"fmt"
	"log/slog"
)

func PanicOnErr(err error) {
	if err != nil {
		fmt.Println("Panicing!!")
		panic(err)
	}
}

func FatalOnErr(err error) {
	if err != nil {
		slog.Error("Error: %s", err)
	}
}
