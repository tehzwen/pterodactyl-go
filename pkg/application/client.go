package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tehzwen/pterodactyl-go/pkg/application/node"
	"github.com/tehzwen/pterodactyl-go/pkg/application/server"
	"github.com/tehzwen/pterodactyl-go/pkg/application/user"
)

type Client struct {
	baseUrl string
	client  *http.Client
	apiKey  string
	Servers *server.ServerApi
	Nodes   *node.NodeApi
	Users   *user.UserApi
}

func (ac *Client) buildHeaders() http.Header {
	return http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ac.apiKey)},
		"Accept":        []string{"Application/vnd.pterodactyl.v1+json"},
		"Content-Type":  []string{"application/json"},
	}
}
func WithApiKey(key string) func(ac *Client) {
	return func(ac *Client) {
		ac.apiKey = key
	}
}

func NewApplicationApi(baseUrl string, opts ...func(a *Client)) (*Client, error) {
	a := &Client{
		client:  &http.Client{Timeout: time.Second * 5},
		baseUrl: baseUrl,
	}

	for _, o := range opts {
		o(a)
	}

	a.Servers = server.NewServerApi(fmt.Sprintf("%s/api/application/servers", a.baseUrl), a.client, a.buildHeaders())
	a.Nodes = node.NewNodeApi(fmt.Sprintf("%s/api/application/nodes", a.baseUrl), a.client, a.buildHeaders())
	a.Users = user.NewUserApi(fmt.Sprintf("%s/api/application/users", a.baseUrl), a.client, a.buildHeaders())
	return a, nil
}
