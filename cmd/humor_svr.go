package main

import (
	"fmt"
	"github.com/wilenceyao/humor-svr/internal"
	"os"
)

func main() {
	err := internal.RunServer()
	if err != nil {
		fmt.Printf("run server err: %+v\n", err)
		os.Exit(1)
	}
}
