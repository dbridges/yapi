package main

import (
	"fmt"
	"os"

	"github.com/dbridges/yapi/app"
)

func main() {
	err := app.New().Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\ntry `yapi --help` for more info\n", err.Error())
	}
}
