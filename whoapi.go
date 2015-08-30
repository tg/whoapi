// Package whoapi allows for communicating with whoAPI.
// Reference: https://whoapi.com/api-documentation.html
package whoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// StatusError is an error received by WhoAPI in JSON response.
type StatusError struct {
	Code Int64  `json:"status"`      // status code
	Desc string `json:"status_desc"` // status description
}

// Error returns a formatted error.
func (e StatusError) Error() string {
	return fmt.Sprintf("WhoAPI: [%d] %s", e.Code, e.Desc)
}

// TODO: add all status code errors
var (
	// ErrTLDDoesNotExist is issued when WhoAPi doesn't recognise TLD.
	ErrTLDDoesNotExist = StatusError{4, "TLD does not exist"}

	// ErrWhoisNotYetSupported is issued when whois server for input value is
	// not yet supported. Note that for whois requests this can still contain
	// a valid and meaningful whois response;
	// see: https://whoapi.com/forum/thick-and-thin-whois-t37.
	ErrWhoisNotYetSupported = StatusError{7, "whois server not yet supported"}

	// ErrInvalidAPIAccount indicates invalid API key.
	ErrInvalidAPIAccount = StatusError{12, "invalid API account"}

	// ErrTooManyRequests is issued when request rate was exceeded.
	ErrTooManyRequests = StatusError{18, "too many requests"}
)

// Int64 is an integer which can be unmarshaled from both number or string
// literal (representing a valid number). This is needed as WhoAPI is
// inconsistent with types, for example sometimes returning "0" (string)
// status code and sometimes 0 (number). If you ever want to parse raw JSON
// response data into a structure, you might consider using this type
// instead of raw int.
type Int64 int64

// UnmarshalJSON decodes Int from raw JSON value.
// If data is a string it removes the quotes trying to decode the value as a raw integer.
func (i *Int64) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	var n int64
	err := json.Unmarshal(data, &n)
	*i = Int64(n)
	return err
}

// Client which talks to WhoAPI. If you have a WHOAPI_KEY environmental
// variable set with a valid API key, you can use this structure uninitialised.
type Client struct {
	// API key to be used for queries. If empty it will be taken from WHOAPI_KEY
	// environmental variable.
	Key string

	// HTTP client to be used for queries.
	// If empty a default client will be used.
	Client *http.Client
}

// Get makes an API request to WhoAPI. Returns raw JSON data and/or and error.
// If a valid response is received, error will be of StatusError type in which
// case also a valid JSON will be returned.
func (c *Client) Get(req string, domain string) (data []byte, err error) {
	data, err = c.GetRaw(req, domain)
	if err != nil {
		return
	}

	// Process response status
	var status StatusError
	err = json.Unmarshal(data, &status)
	if err != nil {
		return
	}

	// Convert status to known error or create a new one
	switch status.Code {
	case 0: // OK
	case ErrTLDDoesNotExist.Code:
		err = ErrTLDDoesNotExist
	case ErrWhoisNotYetSupported.Code:
		err = ErrWhoisNotYetSupported
	case ErrInvalidAPIAccount.Code:
		err = ErrInvalidAPIAccount
	case ErrTooManyRequests.Code:
		err = ErrTooManyRequests
	default:
		err = status
	}

	return
}

// GetRaw makes an API request and returns data without any processing.
func (c *Client) GetRaw(req string, domain string) (data []byte, err error) {
	// Create a default client if not provided
	if c.Client == nil {
		c.Client = &http.Client{Timeout: time.Second * 30}
	}
	// Get api key from environment if not set
	if c.Key == "" {
		c.Key = os.Getenv("WHOAPI_KEY")
	}

	// Prepare url
	params := url.Values{}
	params.Set("apikey", c.Key)
	params.Set("r", req)
	params.Set("domain", domain)
	url := fmt.Sprintf("http://api.whoapi.com?%s", params.Encode())

	// Query API
	res, err := c.Client.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}

	return ioutil.ReadAll(res.Body)
}
