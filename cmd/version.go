package cmd

import (
	"fmt"
)

// VersionString is set at build time
var VersionString string

// Version prints the current version
func Version() error {
	fmt.Printf("yapi version %s\n", VersionString)
	return nil
}
