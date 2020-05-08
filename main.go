package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Version is set at build time
var Version string

var helpFlag bool
var versionFlag bool
var listFlag bool
var nameFlag string

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Display help")
	flag.BoolVar(&helpFlag, "h", false, "Display help")
	flag.BoolVar(&versionFlag, "version", false, "Display version")
	flag.BoolVar(&listFlag, "list", false, "List available route names")
	flag.StringVar(&nameFlag, "name", "", "Run `route_name`")
}

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

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if helpFlag {
		usage()
		return
	}
	if versionFlag {
		version()
		return
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "no source file provided")
		usage()
		return
	}
	fname, line, err := parseFileArg(args[0])
	Must(err)
	if len(fname) <= 0 {
		fmt.Fprintln(os.Stderr, "no source file provided")
		usage()
		return
	}
	cfg, err := NewYAMLConfig(fname)
	Must(err)
	if listFlag {
		list(cfg)
		return
	}
	if len(line) <= 0 && nameFlag == "" {
		fmt.Fprintln(os.Stderr, "no line or route name specified")
		usage()
		return
	}
	if nameFlag == "" {
		lineInt, err := strconv.Atoi(line)
		Must(err)
		nameFlag, err = cfg.FindRequestName(lineInt)
		Must(err)
	}
	c := NewFetchController(cfg)
	err = c.DoRequest(nameFlag)
	Must(err)
}

func list(cfg Config) {
	fmt.Println("Named routes:")
	for _, r := range cfg.RequestNames() {
		fmt.Printf("    %s\n", r)
	}
}

func fetch(cfg Config, name string) {

}

func version() {
	fmt.Printf("yapi version %s\n", Version)
}

func usage() {
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

func parseFileArg(arg string) (filename string, lineNumber string, err error) {
	re := regexp.MustCompile(`(.*?):?(\d+)?$`)
	match := re.FindStringSubmatch(arg)
	if len(match) == 3 {
		return match[1], match[2], nil
	}
	if len(match) == 2 {
		return match[1], "", nil
	}
	return "", "", fmt.Errorf("could not parse '%s'", arg)
}
