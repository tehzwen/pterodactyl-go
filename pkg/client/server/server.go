package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tehzwen/pterodactyl-go/pkg/client/types"
)

type ServerApi struct {
	baseUrl    string
	client     *http.Client
	authHeader http.Header
}

func NewServerApi(baseUrl string, client *http.Client, authHeader http.Header) *ServerApi {
	return &ServerApi{
		baseUrl:    baseUrl,
		client:     client,
		authHeader: authHeader,
	}
}

func (sa *ServerApi) ListServers(ctx context.Context, params types.ListServersParams) ([]types.Server, error) {
	headers := sa.authHeader
	url, err := url.Parse(sa.baseUrl)
	if err != nil {
		return nil, err
	}

	queryParams := url.Query()
	if params.Page > 0 {
		queryParams.Add("page", strconv.Itoa(params.Page))
	}
	if params.PerPage > 0 {
		queryParams.Add("per_page", strconv.Itoa(params.PerPage))
	}
	if params.Include != "" {
		queryParams.Add("include", params.Include)
	}

	var servers []types.Server = make([]types.Server, 0)

	for {
		url.RawQuery = queryParams.Encode()
		fmt.Println(url)
		req := &http.Request{
			Method: "GET",
			Header: headers,
			URL:    url,
		}
		req.WithContext(ctx)

		resp, err := sa.client.Do(req)
		if err != nil {
			return nil, err
		}

		if err := types.CheckResponse(resp); err != nil {
			return nil, err
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var response types.ListServersResponse
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}

		servers = append(servers, response.Servers...)
		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		// increment the url params
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}

	return servers, nil
}
