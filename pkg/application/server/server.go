package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
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
		authHeader: authHeader.Clone(),
	}
}

// SERVERS
func (sa *ServerApi) ListServers(ctx context.Context, request *types.PteroListRequest) ([]types.Server, error) {
	headers := sa.authHeader.Clone()
	url, err := url.Parse(sa.baseUrl)
	if err != nil {
		return nil, err
	}

	queryParams := request.BuildQueryParams(url)

	var servers []types.Server
	for {
		url.RawQuery = queryParams.Encode()
		req := &http.Request{
			Method: "GET",
			Header: headers,
			URL:    url,
		}
		req = req.WithContext(ctx)

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

// handles both internal server requests & get server by external ID requests
func (sa *ServerApi) GetServerDetails(ctx context.Context, serverId int, include []types.GetServerDetailsIncludeField, external bool) (*types.Server, error) {
	headers := sa.authHeader.Clone()
	baseUrl := fmt.Sprintf("%s/%s", sa.baseUrl, strconv.Itoa(serverId))
	if external {
		baseUrl = fmt.Sprintf("%s/external/%s", sa.baseUrl, strconv.Itoa(serverId))
	}

	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()

	if len(include) > 0 {
		includeFields := []string{}
		for _, i := range include {
			includeFields = append(includeFields, string(i))
		}
		queryParams.Add("include", strings.Join(includeFields, ","))
	}

	url.RawQuery = queryParams.Encode()
	req := &http.Request{
		Method: "GET",
		Header: headers,
		URL:    url,
	}

	req = req.WithContext(ctx)
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

	fmt.Printf("bytes - %s\n", string(b))

	// TODO - marshal to server and return

	return nil, nil
}

func (sa *ServerApi) CreateServer(ctx context.Context, request types.CreateServerRequest) (*types.Server, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	headers := sa.authHeader.Clone()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", sa.baseUrl, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header = headers
	req = req.WithContext(ctx)

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

	var server types.Server
	if err := json.Unmarshal(b, &server); err != nil {
		return nil, err
	}

	return &server, nil
}

func (sa *ServerApi) updateServer(ctx context.Context, path string, request types.UpdateServerRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := sa.authHeader.Clone()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", path, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	req.Header = headers
	req = req.WithContext(ctx)

	resp, err := sa.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

func (sa *ServerApi) UpdateServerDetails(ctx context.Context, request types.UpdateServerDetailsRequest) error {
	return sa.updateServer(ctx, fmt.Sprintf("%s/%d/details", sa.baseUrl, request.ServerId), request)
}

func (sa *ServerApi) UpdateServerBuild(ctx context.Context, request types.UpdateServerBuildRequest) error {
	return sa.updateServer(ctx, fmt.Sprintf("%s/%d/build", sa.baseUrl, request.ServerId), request)
}

func (sa *ServerApi) UpdateServerStartup(ctx context.Context, request types.UpdateServerStartupRequest) error {
	return sa.updateServer(ctx, fmt.Sprintf("%s/%d/startup", sa.baseUrl, request.ServerId), request)
}

// helper for all POST endpoints
func (sa *ServerApi) postServer(ctx context.Context, path string) error {
	headers := sa.authHeader.Clone()
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}

	req.Header = headers
	req = req.WithContext(ctx)

	resp, err := sa.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

func (sa *ServerApi) SuspendServer(ctx context.Context, serverId int) error {
	return sa.postServer(ctx, fmt.Sprintf("%s/%d/suspend", sa.baseUrl, serverId))
}

func (sa *ServerApi) UnsuspendServer(ctx context.Context, serverId int) error {
	return sa.postServer(ctx, fmt.Sprintf("%s/%d/unsuspend", sa.baseUrl, serverId))
}

func (sa *ServerApi) ReinstallServer(ctx context.Context, serverId int) error {
	return sa.postServer(ctx, fmt.Sprintf("%s/%d/reinstall", sa.baseUrl, serverId))
}

func (sa *ServerApi) DeleteServer(ctx context.Context, serverId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", sa.baseUrl, serverId), nil)
	if err != nil {
		return err
	}

	req.Header = sa.authHeader
	req = req.WithContext(ctx)

	resp, err := sa.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}
