package main

import (
	"fmt"
	"os"
)

// Must exits the program in the error is not nil
func Must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
