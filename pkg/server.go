package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type GetServerDetailsIncludeField string

var (
	ServerDetailsAllocations GetServerDetailsIncludeField = "allocations"
	ServerDetailsUser        GetServerDetailsIncludeField = "user"
	ServerDetailsSubUsers    GetServerDetailsIncludeField = "subusers"
	ServerDetailsPack        GetServerDetailsIncludeField = "pack"
	ServerDetailsNest        GetServerDetailsIncludeField = "nest"
	ServerDetailsEgg         GetServerDetailsIncludeField = "egg"
	ServerDetailsVariables   GetServerDetailsIncludeField = "variables"
	ServerDetailsLocation    GetServerDetailsIncludeField = "location"
	ServerDetailsNode        GetServerDetailsIncludeField = "node"
	ServerDetailsDatabases   GetServerDetailsIncludeField = "databases"
	ServerDetailsBackups     GetServerDetailsIncludeField = "backups"
)

// handles both internal server requests & get server by external ID requests
func (ac *ApplicationClient) GetServerDetails(ctx context.Context, serverId int, include []GetServerDetailsIncludeField, external bool) (*Server, error) {
	headers := ac.buildHeaders()
	baseUrl := fmt.Sprintf("%s/api/application/servers/%s", ac.baseUrl, strconv.Itoa(serverId))
	if external {
		baseUrl = fmt.Sprintf("%s/api/application/servers/external/%s", ac.baseUrl, strconv.Itoa(serverId))
	}

	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()

	if include != nil && len(include) > 0 {
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

	fmt.Printf("bytes - %s\n", string(b))

	// TODO - marshal to server and return

	return nil, nil
}

func (ac *ApplicationClient) CreateServer(ctx context.Context, request CreateServerRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.buildHeaders()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/application/servers", ac.baseUrl), bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	req.Header = headers
	req.WithContext(ctx)

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("bytes - %s\n", string(b))

	return nil
}

type UpdateServerRequest interface {
	Validate() error
}

func (ac *ApplicationClient) updateServer(ctx context.Context, path string, request UpdateServerRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.buildHeaders()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", path, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	req.Header = headers
	req.WithContext(ctx)

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("bytes - %s\n", string(b))

	return nil
}

type UpdateServerDetailsRequest struct {
	ServerId    string
	Name        string `json:"name"`
	UserId      int    `json:"user"`
	ExternalId  string `json:"external_id"`
	Description string `json:"description"`
}

func (usr UpdateServerDetailsRequest) Validate() error {
	return nil
}

func (ac *ApplicationClient) UpdateServerDetails(ctx context.Context, request UpdateServerDetailsRequest) error {
	return ac.updateServer(ctx, fmt.Sprintf("%s/api/application/servers/%s/details", ac.baseUrl, request.ServerId), request)
}

type UpdateServerBuildRequest struct {
	ServerId          int
	Allocation        int           `json:"allocation"`
	Memory            int           `json:"memory"`
	Swap              int           `json:"swap"`
	Disk              int           `json:"disk"`
	IO                int           `json:"io"`
	CPU               int           `json:"cpu"`
	Threads           string        `json:"threads,omitempty"`
	FeatureLimits     FeatureLimits `json:"feature_limits"`
	AddAllocations    []int         `json:"add_allocations,omitempty"`
	RemoveAllocations []int         `json:"remove_allocations,omitempty"`
	OOMDisabled       bool          `json:"oom_disabled,omitempty"`
}

func (r UpdateServerBuildRequest) Validate() error {
	if r.Allocation == 0 {
		return errors.New("allocation is required")
	}
	if r.Memory == 0 {
		return errors.New("memory is required")
	}
	if r.Disk == 0 {
		return errors.New("disk is required")
	}
	if r.IO == 0 {
		return errors.New("io is required")
	}
	if r.CPU == 0 {
		return errors.New("cpu is required")
	}

	return nil
}

func (ac *ApplicationClient) UpdateServerBuild(ctx context.Context, request UpdateServerBuildRequest) error {
	return ac.updateServer(ctx, fmt.Sprintf("%s/api/application/servers/%s/build", ac.baseUrl, request.ServerId), request)
}

type UpdateServerStartupRequest struct {
	ServerId    string
	Startup     string            `json:"startup"`
	Environment map[string]string `json:"environment"`
	Egg         int               `json:"egg"`
	Image       string            `json:"image,omitempty"`
	SkipScripts bool              `json:"skip_scripts,omitempty"`
}

func (r UpdateServerStartupRequest) Validate() error {
	if r.Startup == "" {
		return errors.New("startup is required")
	}
	if r.Environment == nil {
		return errors.New("environment is required")
	}
	if r.Egg == 0 {
		return errors.New("egg is required")
	}

	return nil
}

func (ac *ApplicationClient) UpdateServerStartup(ctx context.Context, request UpdateServerStartupRequest) error {
	return ac.updateServer(ctx, fmt.Sprintf("%s/api/application/servers/%s/build", ac.baseUrl, request.ServerId), request)
}
