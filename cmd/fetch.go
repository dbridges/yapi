package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dbridges/yapi/client"
	"github.com/dbridges/yapi/config"
)

type fetch struct {
	cfg config.Config
}

func Fetch(cfg config.Config, name string) error {
	f := &fetch{cfg: cfg}
	return f.Run(name)
}

func (f *fetch) Run(name string) error {
	r, err := f.cfg.NewRequest(name)
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
	client, err := client.New(f.cfg.SessionName())
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return f.printResponse(resp)
}

func (f *fetch) printResponse(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println("")

	// Print Headers if wanted
	if f.cfg.DisplayHeaders() {
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
		if err != nil {
			return err
		}
		fmt.Println(output.String())
	} else {
		fmt.Println(string(body))
	}
	return nil
}
