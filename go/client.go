package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

type Client struct {
	apiKey       string
	Organization string
}

// NewClient creates a new client
func NewClient(apiKey string, organization string) *Client {
	return &Client{
		apiKey:       apiKey,
		Organization: organization,
	}
}

// Post makes a post request
func (c *Client) Post(url string, input any) (response []byte, err error) {
	response = make([]byte, 0)

	rJson, err := json.Marshal(input)
	if err != nil {
		return response, err
	}

	resp, err := c.Call(http.MethodPost, url, bytes.NewReader(rJson))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	response, err = io.ReadAll(resp.Body)
	return response, err
}

// Get makes a get request
func (c *Client) Get(url string, input any) (response []byte, err error) {
	if input != nil {
		vals, _ := query.Values(input)
		query := vals.Encode()

		if query != "" {
			url += "?" + query
		}
	}

	resp, err := c.Call(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	response, err = io.ReadAll(resp.Body)
	return response, err
}

// Call makes a request
func (c *Client) Call(method string, url string, body io.Reader) (response *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return response, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	if c.Organization != "" {
		req.Header.Add("OpenAI-Organization", c.Organization)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
