package errors

import "fmt"

func PanicOnErr(err error) {
	if err != nil {
		fmt.Println("Panicing!!")
		panic(err)
	}
}
