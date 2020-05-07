package main

import (
	"net/http"
	"os"
	"path"
	"time"

	cookiejar "github.com/juju/persistent-cookiejar"
	"golang.org/x/net/publicsuffix"
)

// Client manages fetching requests and cookies
type Client struct {
	client      *http.Client
	sessionName string
}

// NewClient creates a new client with the provided session Name
func NewClient(sessionName string) (*Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	savePath := path.Join(home, ".cache/yapi", sessionName+".jar.json")
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
		Filename:         savePath,
	})
	if err != nil {
		return nil, err
	}
	c := &Client{
		client: &http.Client{
			Timeout: time.Second * 10,
			Jar:     jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		sessionName: sessionName,
	}
	return c, nil
}

// Do performs the given request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.sessionName != "" {
		jar, ok := c.client.Jar.(*cookiejar.Jar)
		if ok {
			jar.Save()
		}
	}
	return resp, nil
}
