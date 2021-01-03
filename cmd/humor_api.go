package main

import (
	"fmt"
	"github.com/wilenceyao/humor-api/internal"
	"os"
)

func main() {
	err := internal.RunServer()
	if err != nil {
		fmt.Println(fmt.Sprintf("run server err: %+v", err))
		os.Exit(1)
	}
}
