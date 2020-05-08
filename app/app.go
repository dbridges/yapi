package app

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"

	"github.com/dbridges/yapi/cmd"
	"github.com/dbridges/yapi/config"
)

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

// App is a container to run common commands
type App struct{}

// New returns a new app instance
func New() *App {
	return &App{}
}

// Run is the main entry point to the app
func (app *App) Run() error {
	flag.Usage = cmd.Usage
	flag.Parse()
	args := flag.Args()

	// Easy flags
	if helpFlag {
		return cmd.Help()
	}
	if versionFlag {
		return cmd.Version()
	}

	// Ensure we have a proper source file
	if len(args) != 1 {
		return fmt.Errorf("no source file provided")
	}
	fname, line, err := parseFileArg(args[0])
	if err != nil {
		return err
	}
	if len(fname) <= 0 {
		return fmt.Errorf("no source file provided")
	}

	// Load the config source file
	cfg, err := config.NewYAMLConfig(fname)
	if err != nil {
		return err
	}
	if listFlag {
		return cmd.List(cfg)
	}
	if len(line) <= 0 && nameFlag == "" {
		return fmt.Errorf("no line or route name specified")
	}
	if nameFlag == "" {
		lineInt, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		nameFlag, err = cfg.FindRequestName(lineInt)
		if err != nil {
			return err
		}
	}
	return cmd.Fetch(cfg, nameFlag)
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
