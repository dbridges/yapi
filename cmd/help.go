package cmd

import (
	"flag"
	"fmt"
	"os"
)

var usageHeader = `
EXAMPLES
  Run a route:
      yapi --name=route path/to/my.yapi.yml
      yapi path/to/my.yapi.yml:10
  
  List available routes:
      yapi --list path/to/my.yapi.yml

ARGUMENTS
  source_file
	A file path with optional line number. The route nearest to line number will be
	run. If no line number is given the --name flag is required.

OPTIONS
`

var usageFooter = `
FURTHER INFORMATION
  Visit https://github.com/dbridges/yapi for more information
`

// Help prints the help text
func Help() error {
	Usage()
	return nil
}

// Usage prints the help text
func Usage() {
	fmt.Fprintln(os.Stderr, "Usage: yapi [opts] source_file")
	fmt.Fprintf(os.Stderr, usageHeader)
	flag.VisitAll(func(f *flag.Flag) {
		argName, usage := flag.UnquoteUsage(f)
		if len(f.Name) == 1 {
			fmt.Fprintf(os.Stderr, "  -%s %s\n", f.Name, argName)
		} else {
			fmt.Fprintf(os.Stderr, "  --%s %s\n", f.Name, argName)
		}
		fmt.Fprintf(os.Stderr, "      %s\n", usage)
	})
	fmt.Fprintf(os.Stderr, usageFooter)
}
