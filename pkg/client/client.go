// client represents the client API of pterodactyl
// see https://pteroapi.com/docs/api/client for more info
package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tehzwen/pterodactyl-go/pkg/client/server"
)

type Client struct {
	baseUrl string
	client  *http.Client
	apiKey  string

	Servers *server.ServerApi
}

func WithApiKey(key string) func(ac *Client) {
	return func(ac *Client) {
		ac.apiKey = key
	}
}

func (ac *Client) buildHeaders() http.Header {
	return http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ac.apiKey)},
		"Accept":        []string{"Application/vnd.pterodactyl.v1+json"},
		"Content-Type":  []string{"application/json"},
	}
}

func NewClientApi(baseUrl string, opts ...func(a *Client)) (*Client, error) {
	c := &Client{
		baseUrl: baseUrl,
		client:  &http.Client{Timeout: time.Second * 5},
	}

	for _, o := range opts {
		o(c)
	}

	c.Servers = server.NewServerApi(fmt.Sprintf("%s/api/client", baseUrl), c.client, c.buildHeaders())
	return c, nil
}
