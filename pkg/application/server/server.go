package server

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
		authHeader: authHeader,
	}
}

type APIErrorResponse struct {
	Errors []struct {
		Code   string `json:"code"`
		Status string `json:"status"`
		Detail string `json:"detail"`
		Source struct {
			Field string `json:"field"`
		} `json:"source"`
	} `json:"errors"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.Message)
}

func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	var errResp APIErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && len(errResp.Errors) > 0 {
		var msgs []string
		for _, e := range errResp.Errors {
			if e.Source.Field != "" {
				msgs = append(msgs, fmt.Sprintf("%s (%s: %s)", e.Detail, e.Source.Field, e.Code))
			} else {
				msgs = append(msgs, fmt.Sprintf("%s (%s)", e.Detail, e.Code))
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.Join(msgs, "; "),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    fallbackStatusMessage(resp.StatusCode),
	}
}

func fallbackStatusMessage(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "invalid input data"
	case http.StatusUnauthorized:
		return "invalid API key"
	case http.StatusForbidden:
		return "insufficient permissions"
	case http.StatusNotFound:
		return "server does not exist"
	case http.StatusUnprocessableEntity:
		return "invalid field values"
	case http.StatusTooManyRequests:
		return "rate limit exceeded"
	default:
		return http.StatusText(code)
	}
}

// SERVERS
func (ac *ServerApi) ListServers(ctx context.Context, request *types.PteroListRequest) ([]types.Server, error) {
	headers := ac.authHeader
	url, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, "/api/application/servers"))
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
		req.WithContext(ctx)

		resp, err := ac.client.Do(req)
		if err != nil {
			return nil, err
		}

		if err := CheckResponse(resp); err != nil {
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
func (ac *ServerApi) GetServerDetails(ctx context.Context, serverId int, include []GetServerDetailsIncludeField, external bool) (*types.Server, error) {
	headers := ac.authHeader
	baseUrl := fmt.Sprintf("%s/api/application/servers/%s", ac.baseUrl, strconv.Itoa(serverId))
	if external {
		baseUrl = fmt.Sprintf("%s/api/application/servers/external/%s", ac.baseUrl, strconv.Itoa(serverId))
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

	req.WithContext(ctx)
	resp, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
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

func (ac *ServerApi) CreateServer(ctx context.Context, request CreateServerRequest) (*types.Server, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	headers := ac.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/application/servers", ac.baseUrl), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header = headers
	req.WithContext(ctx)

	resp, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
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

type UpdateServerRequest interface {
	Validate() error
}

func (ac *ServerApi) updateServer(ctx context.Context, path string, request UpdateServerRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.authHeader
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

	if err := CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

type UpdateServerDetailsRequest struct {
	ServerId    int
	Name        string `json:"name"`
	UserId      int    `json:"user"`
	ExternalId  string `json:"external_id"`
	Description string `json:"description"`
}

func (usr UpdateServerDetailsRequest) Validate() error {
	if usr.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (ac *ServerApi) UpdateServerDetails(ctx context.Context, request UpdateServerDetailsRequest) error {
	return ac.updateServer(ctx, fmt.Sprintf("%s/api/application/servers/%d/details", ac.baseUrl, request.ServerId), request)
}

type UpdateServerBuildRequest struct {
	ServerId          int
	Allocation        int                 `json:"allocation"`
	Memory            int                 `json:"memory"`
	Swap              int                 `json:"swap"`
	Disk              int                 `json:"disk"`
	IO                int                 `json:"io"`
	CPU               int                 `json:"cpu"`
	Threads           string              `json:"threads,omitempty"`
	FeatureLimits     types.FeatureLimits `json:"feature_limits"`
	AddAllocations    []int               `json:"add_allocations,omitempty"`
	RemoveAllocations []int               `json:"remove_allocations,omitempty"`
	OOMDisabled       bool                `json:"oom_disabled,omitempty"`
}

func (r UpdateServerBuildRequest) Validate() error {
	if r.Allocation == 0 {
		return errors.New("allocation is required")
	}
	if r.Memory == 0 {
		return errors.New("memory is required")
	}
	if r.Disk < 0 {
		return errors.New("disk is required")
	}
	if r.IO == 0 {
		return errors.New("io is required")
	}
	if r.CPU < 0 {
		return errors.New("cpu is required")
	}

	return nil
}

func (ac *ServerApi) UpdateServerBuild(ctx context.Context, request UpdateServerBuildRequest) error {
	return ac.updateServer(ctx, fmt.Sprintf("%s/api/application/servers/%d/build", ac.baseUrl, request.ServerId), request)
}

type UpdateServerStartupRequest struct {
	ServerId    int
	Startup     string         `json:"startup"`
	Environment map[string]any `json:"environment"`
	Egg         int            `json:"egg"`
	Image       string         `json:"image,omitempty"`
	SkipScripts bool           `json:"skip_scripts"`
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

func (ac *ServerApi) UpdateServerStartup(ctx context.Context, request UpdateServerStartupRequest) error {
	return ac.updateServer(ctx, fmt.Sprintf("%s/api/application/servers/%d/startup", ac.baseUrl, request.ServerId), request)
}

// helper for all POST endpoints
func (ac *ServerApi) postServer(ctx context.Context, path string) error {
	headers := ac.authHeader
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}

	req.Header = headers
	req.WithContext(ctx)

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}

	if err := CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

func (ac *ServerApi) SuspendServer(ctx context.Context, serverId int) error {
	return ac.postServer(ctx, fmt.Sprintf("%s/api/application/servers/%d/suspend", ac.baseUrl, serverId))
}

func (ac *ServerApi) UnsuspendServer(ctx context.Context, serverId int) error {
	return ac.postServer(ctx, fmt.Sprintf("%s/api/application/servers/%d/unsuspend", ac.baseUrl, serverId))
}

func (ac *ServerApi) ReinstallServer(ctx context.Context, serverId int) error {
	return ac.postServer(ctx, fmt.Sprintf("%s/api/application/servers/%d/reinstall", ac.baseUrl, serverId))
}

func (ac *ServerApi) DeleteServer(ctx context.Context, serverId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/application/servers/%d", ac.baseUrl, serverId), nil)
	if err != nil {
		return err
	}

	req.Header = ac.authHeader
	req.WithContext(ctx)

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}

	if err := CheckResponse(resp); err != nil {
		return err
	}

	return nil
}
