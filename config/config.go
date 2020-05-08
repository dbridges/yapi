package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// Config is an interface to a given YAPI configuration
type Config interface {
	DisplayHeaders() bool
	RequestNames() []string
	FindRequestName(int) (string, error)
	NewRequest(string) (*Request, error)
	SessionName() string
}

// Request stores the details for a single request
type Request struct {
	Path    string
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    string
}

// YAMLConfig implements Config backed by a YAML file
type YAMLConfig struct {
	Root     string
	Headers  map[string]string
	Output   map[string]bool
	Session  string
	Requests map[string]interface{}
	data     []byte
}

// TODO: dynamically generate this using reflection
// YAMLConfigSettingsKeys are reserved for settings use
var YAMLConfigSettingsKeys = map[string]bool{
	"root":    true,
	"headers": true,
	"output":  true,
	"session": true,
}

// NewYAMLConfig loads a config for a given filename
func NewYAMLConfig(fname string) (*YAMLConfig, error) {
	cfg := YAMLConfig{}
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	cfg.data = data
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg.Requests)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// NewRequest creates a new request based on the current config
func (cfg *YAMLConfig) NewRequest(name string) (*Request, error) {
	var r Request
	cfgReq, ok := cfg.Requests[name]
	if !ok {
		return nil, fmt.Errorf("unable to find request named `%s`", name)
	}
	err := mapstructure.Decode(cfgReq, &r)
	if err != nil {
		return nil, err
	}
	if r.Method == "" {
		r.Method = "GET"
	}
	r.URL = cfg.Root + r.Path
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	for k, v := range cfg.Headers {
		if _, ok := r.Headers[k]; !ok {
			r.Headers[k] = v
		}
	}
	return &r, nil
}

// RequestNames lists available request names
func (cfg *YAMLConfig) RequestNames() []string {
	items := make([]string, 0, len(cfg.Requests))
	for k := range cfg.Requests {
		if !YAMLConfigSettingsKeys[k] {
			items = append(items, k)
		}
	}
	return items
}

// FindRequestName tries to find the nearest (going upward) request for a
// given line number
func (cfg *YAMLConfig) FindRequestName(line int) (string, error) {
	re := regexp.MustCompile(`(^\S.*):$`)
	lines := strings.Split(string(cfg.data), "\n")
	if line >= len(lines) {
		return "", fmt.Errorf("file has no line %d", line)
	}
	row := line
	for row > 0 {
		match := re.FindStringSubmatch(lines[row])
		if len(match) == 2 && !YAMLConfigSettingsKeys[match[1]] {
			return match[1], nil
		}
		row--
	}
	return "", fmt.Errorf("unable to find request near line %d", line)
}

// DisplayHeaders returns whether or not the config wants headers to be
// displayed in the output.
func (cfg *YAMLConfig) DisplayHeaders() bool {
	val, ok := cfg.Output["headers"]
	return val && ok
}

// SessionName returns the name used for cookie storage
func (cfg *YAMLConfig) SessionName() string {
	return cfg.Session
}
