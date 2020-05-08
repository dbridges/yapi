package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// FetchController manages performing a request and displaying its response.
type FetchController struct {
	cfg Config
}

// NewFetchController creates a new FetchController
func NewFetchController(cfg Config) *FetchController {
	return &FetchController{
		cfg: cfg,
	}
}

// DoRequest performs the given request and returns an http response
func (c *FetchController) DoRequest(name string) error {
	r, err := c.cfg.NewRequest(nameFlag)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return err
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
	client, err := NewClient(c.cfg.SessionName())
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	c.PrintResponse(resp)
	return nil
}

// PrintResponse displays the response according to the configuration
func (c *FetchController) PrintResponse(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	Must(err)
	fmt.Println(resp.Status)
	fmt.Println("")

	// Print Headers if wanted
	if c.cfg.DisplayHeaders() {
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
