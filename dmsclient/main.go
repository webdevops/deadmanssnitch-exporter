package dmsclient

import (
	"errors"
	"fmt"
	"sync/atomic"

	"gopkg.in/resty.v1"
)

const (
	DmsApiUrl = "https://api.deadmanssnitch.com/v1/"
)

type Client struct {
	// RequestCount has to be the first words
	// in order to be 64-aligned on 32-bit architectures.
	RequestCount   uint64
	RequestRetries int

	accessToken *string

	restClient *resty.Client
}

func NewClient(token string) *Client {
	c := Client{}
	c.SetAccessToken(token)
	c.Init()

	return &c
}

func (c *Client) Init() {
	c.RequestCount = 0
	c.SetRetries(3)
}

func (c *Client) SetRetries(v int) {
	c.RequestRetries = v

	if c.restClient != nil {
		c.restClient.SetRetryCount(c.RequestRetries)
	}
}

func (c *Client) SetUserAgent(v string) {
	c.rest().SetHeader("User-Agent", v)
}

func (c *Client) SetAccessToken(token string) {
	c.accessToken = &token
	c.rest().SetBasicAuth(*c.accessToken, "")
}

func (c *Client) rest() *resty.Client {
	if c.restClient == nil {
		c.restClient = resty.New()
		c.restClient.SetHostURL(DmsApiUrl)
		c.restClient.SetHeader("Accept", "application/json")
		c.restClient.SetBasicAuth(*c.accessToken, "")
		c.restClient.SetRetryCount(c.RequestRetries)
		c.restClient.OnBeforeRequest(c.restOnBeforeRequest)
		c.restClient.OnAfterResponse(c.restOnAfterResponse)
	}

	return c.restClient
}

func (c *Client) restOnBeforeRequest(client *resty.Client, request *resty.Request) (err error) {
	atomic.AddUint64(&c.RequestCount, 1)
	return
}

func (c *Client) restOnAfterResponse(client *resty.Client, response *resty.Response) (err error) {
	return
}

func (c *Client) GetRequestCount() float64 {
	requestCount := atomic.LoadUint64(&c.RequestCount)
	return float64(requestCount)
}

func (c *Client) checkResponse(response *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if response != nil {
		// check status code
		statusCode := response.StatusCode()
		if statusCode != 200 {
			return fmt.Errorf("response status code is %v (expected 200), url: %v", statusCode, response.Request.URL)
		}
	} else {
		return errors.New("response is nil")
	}

	return nil
}
