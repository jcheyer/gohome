package gohome

import (
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	authToken     string
	host          string
	pingOnStartup bool
	resty         *resty.Client
}

type Option func(*Client)

func New(options ...Option) (*Client, error) {
	c := &Client{
		host:          "http://hassio.local:8123",
		pingOnStartup: false,
	}

	for _, option := range options {
		option(c)
	}

	if c.resty == nil {
		c.resty = resty.New()
	}

	c.resty.SetHostURL(c.host)
	if c.authToken != "" {
		c.resty.SetAuthToken(c.authToken)
	}

	if c.pingOnStartup {
		return c, c.Ping()
	}

	return c, nil
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
