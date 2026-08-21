package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// SERVERS
func (ac *ApplicationClient) ListServers(ctx context.Context, request *PteroListRequest) ([]Server, error) {
	headers := ac.buildHeaders()
	url, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, "/api/application/servers"))
	if err != nil {
		return nil, err
	}

	queryParams := request.buildQueryParams(url)

	var servers []Server
	for {
		url.RawQuery = queryParams.Encode()
		req := &http.Request{
			Method: "GET",
			Header: headers,
			URL:    url,
		}
		req.WithContext(ctx)

		resp, err := ac.client.Do(req)
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var response ListServersResponse
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
