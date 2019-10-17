package gohome

import (
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	authToken         string
	host              string
	pingOnStartup     bool
	discoverOnStartup bool
	discoveryInfo     discoveryInfo
	resty             *resty.Client
}

type Option func(*Client)

func New(options ...Option) (*Client, error) {
	c := &Client{
		host:              "http://hassio.local:8123",
		pingOnStartup:     false,
		discoverOnStartup: true,
	}

	for _, option := range options {
		option(c)
	}

	if c.resty == nil {
		c.resty = resty.New()
	}

	c.resty.SetHostURL(c.host).SetHeader("Content-Type", "application/json")

	if c.discoverOnStartup {
		err := c.discover()
		if err != nil {
			return c, err
		}
	}

	if c.authToken != "" {
		c.resty.SetAuthToken(c.authToken)
	}

	if c.discoveryInfo.RequiresAPIPassword && c.authToken == "" {
		return c, errors.New("requires auth but no auth token set")
	}

	if c.pingOnStartup {
		return c, c.Ping()
	}

	return c, nil
}

func NoDiscovery() Option {
	return func(c *Client) {
		c.discoverOnStartup = false
	}
}

func WithClient(hc *http.Client) Option {
	return func(c *Client) {
		c.resty = resty.NewWithClient(hc)
	}
}

func WithHost(h string) Option {
	return func(c *Client) {
		c.host = h
	}
}

func WithPing() Option {
	return func(c *Client) {
		c.pingOnStartup = true
	}
}

func WithAuthToken(token string) Option {
	return func(c *Client) {
		c.authToken = token
	}
}

func (c *Client) Ping() error {
	resp, err := c.resty.R().Get("/api/")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(string(resp.Body()))
	}

	return nil
}

func (c *Client) SetDebug(d bool) {
	c.resty.SetDebug(d)
}

func (c *Client) discover() error {
	resp, err := c.resty.R().
		SetResult(&discoveryInfo{}).
		Get("/api/discovery_info")
	if err != nil {
		return err
	}

	disco := resp.Result().(*discoveryInfo)
	c.discoveryInfo = *disco

	return nil
}

func (c *Client) Version() string {
	var err error
	if c.discoveryInfo.Version == "" {
		err = c.discover()
	}

	if err != nil {
		return ""
	}

	return c.discoveryInfo.Version
}
