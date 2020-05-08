package cmd

import (
	"fmt"

	"github.com/dbridges/yapi/config"
)

func List(cfg config.Config) error {
	fmt.Println("Named routes:")
	for _, r := range cfg.RequestNames() {
		fmt.Printf("    %s\n", r)
	}
	return nil
}
