package gohome

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
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

func (c *Client) Ping() error {
	_, err := c.resty.R().Get("/api")
	/*if err != nil {
		return err
	}*/
	return err
}