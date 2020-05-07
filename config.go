package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// Request stores the details for a single request
type Request struct {
	Path    string
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    string
}

// Config stores the yaml project configuration data
type Config struct {
	Root     string
	Headers  map[string]string
	Output   map[string]bool
	Session  string
	Requests map[string]interface{}
	data     []byte
}

// TODO: dynamically generate this using reflection
var ConfigSettingsKeys = map[string]bool{
	"root":    true,
	"headers": true,
	"output":  true,
	"session": true,
}

// NewConfig loads a config for a given filename
func NewConfig(fname string) (*Config, error) {
	cfg := Config{}
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
func (cfg *Config) NewRequest(name string) (*Request, error) {
	var r Request
	err := mapstructure.Decode(cfg.Requests[name], &r)
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

func (cfg *Config) RequestNames() []string {
	items := make([]string, 0, len(cfg.Requests))
	for k := range cfg.Requests {
		if !ConfigSettingsKeys[k] {
			items = append(items, k)
		}
	}
	return items
}

// RequestNameForLine tries to find the nearest (going upward) request for a
// given line number
func (cfg *Config) RequestNameForLine(line int) (string, error) {
	re := regexp.MustCompile(`(^\S.*):$`)
	lines := strings.Split(string(cfg.data), "\n")
	if line >= len(lines) {
		return "", fmt.Errorf("file has no line %d", line)
	}
	row := line
	for row > 0 {
		match := re.FindStringSubmatch(lines[row])
		if len(match) == 2 && !ConfigSettingsKeys[match[1]] {
			return match[1], nil
		}
		row--
	}
	return "", fmt.Errorf("unable to find request near line %d", line)
}

// PrintResponse displays the response according to the configuration
func (cfg *Config) PrintResponse(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	Must(err)
	fmt.Println(resp.Status)
	fmt.Println("")

	// Print Headers if wanted
	if val, ok := cfg.Output["headers"]; val && ok {
		for k := range resp.Header {
			fmt.Printf("%s: %s\n", k, resp.Header.Get(k))
		}
		fmt.Println("")
	}

	// Print Body
	if strings.Index(resp.Header.Get("Content-Type"), "application/json") >= 0 {
		// Pretty print JSON
		var output bytes.Buffer
		err := json.Indent(&output, body, "", "  ")
		Must(err)
		fmt.Println(output.String())
	} else {
		fmt.Println(string(body))
	}
}

// DoRequest performs the given request and returns a http response
func (cfg *Config) DoRequest(r *Request) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return nil, err
	}
	if len(r.Params) > 0 {
		q := req.URL.Query()
		for k, v := range r.Params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	client, err := NewClient(cfg.Session)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
