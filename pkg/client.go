package pkg

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ApplicationClient struct {
	baseUrl string
	client  *http.Client
	apiKey  string
}

func (pt *PteroListRequest) buildQueryParams(url *url.URL) url.Values {
	queryParams := url.Query()

	if pt.Include != nil {
		queryParams.Add("include", strings.Join(pt.Include, ","))
	}

	if pt.Page != 0 {
		queryParams.Add("page", strconv.Itoa(pt.Page))
	}

	if pt.PerPage != 0 {
		queryParams.Add("per_page", strconv.Itoa(pt.PerPage))
	}
	return queryParams
}

func (ac *ApplicationClient) buildHeaders() http.Header {
	return http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ac.apiKey)},
		"Accept":        []string{"Application/vnd.pterodactyl.v1+json"},
		"Content-Type":  []string{"application/json"},
	}
}
func WithApiKey(key string) func(ac *ApplicationClient) {
	return func(ac *ApplicationClient) {
		ac.apiKey = key
	}
}

func NewApplicationClient(baseUrl string, opts ...func(a *ApplicationClient)) (*ApplicationClient, error) {
	a := &ApplicationClient{
		client:  &http.Client{Timeout: time.Second * 5},
		baseUrl: baseUrl,
	}

	for _, o := range opts {
		o(a)
	}

	return a, nil
}
