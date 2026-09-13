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

func (sa *ServerApi) GetServerDetails(ctx context.Context, serverIdentifier string) (*types.Server, error) {
	headers := sa.authHeader
	url, err := url.Parse(fmt.Sprintf("%s/servers/%s", sa.baseUrl, serverIdentifier))
	if err != nil {
		return nil, err
	}

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

	var response types.Server
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (sa *ServerApi) GetServerResources(ctx context.Context, serverIdentifier string) (*types.ServerResources, error) {
	headers := sa.authHeader
	url, err := url.Parse(fmt.Sprintf("%s/servers/%s/resources", sa.baseUrl, serverIdentifier))
	if err != nil {
		return nil, err
	}

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

	var response types.ServerResources
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (sa *ServerApi) GetConsoleAccess(ctx context.Context, serverIdentifier string) (*types.ConsoleAccessDetails, error) {
	headers := sa.authHeader
	url, err := url.Parse(fmt.Sprintf("%s/servers/%s/websocket", sa.baseUrl, serverIdentifier))
	if err != nil {
		return nil, err
	}

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

	var response types.ConsoleAccessDetails
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (sa *ServerApi) GetServerActivity(ctx context.Context, serverIdentifier string) ([]types.ActivityLog, error) {
	headers := sa.authHeader
	url, err := url.Parse(fmt.Sprintf("%s/servers/%s/activity", sa.baseUrl, serverIdentifier))
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()

	var activityLogs []types.ActivityLog = make([]types.ActivityLog, 0)

	for {
		url.RawQuery = queryParams.Encode()
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

		var response types.GetServerActivityResponse
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}

		activityLogs = append(activityLogs, response.Data...)
		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}

	return activityLogs, nil
}

func (sa *ServerApi) ReinstallServer(ctx context.Context, serverIdentifier string) error {
	headers := sa.authHeader
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/servers/%s/settings/reinstall", sa.baseUrl, serverIdentifier), nil)
	if err != nil {
		return err
	}

	req.Header = headers
	req.WithContext(ctx)

	resp, err := sa.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

func (sa *ServerApi) ManagePower(ctx context.Context, request types.ManagePowerRequest) error {
	headers := sa.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/servers/%s/power", sa.baseUrl, request.ServerIdentifier), bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	req.Header = headers
	req.WithContext(ctx)

	resp, err := sa.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}
