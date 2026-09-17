package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

type UserApi struct {
	baseUrl    string
	client     *http.Client
	authHeader http.Header
}

func NewUserApi(baseUrl string, client *http.Client, authHeader http.Header) *UserApi {
	return &UserApi{
		baseUrl:    baseUrl,
		client:     client,
		authHeader: authHeader.Clone(),
	}
}

func (ua *UserApi) ListUsers(ctx context.Context, filters types.ListUsersFilters) ([]types.User, error) {
	headers := ua.authHeader
	url, err := url.Parse(ua.baseUrl)
	if err != nil {
		return nil, err
	}

	queryParams := url.Query()
	if filters.Page > 0 {
		queryParams.Add("page", strconv.Itoa(filters.Page))
	}
	if filters.PerPage > 0 {
		queryParams.Add("per_page", strconv.Itoa(filters.PerPage))
	}
	if filters.Email != "" {
		queryParams.Add("filter[email]", filters.Email)
	}
	if filters.UUID != "" {
		queryParams.Add("filter[uuid]", filters.UUID)
	}
	if filters.Username != "" {
		queryParams.Add("filter[username]", filters.Username)
	}
	if filters.ExternalID != "" {
		queryParams.Add("filter[external_id]", filters.ExternalID)
	}
	if filters.Sort != "" {
		queryParams.Add("sort", filters.Sort)
	}
	if filters.Include != "" {
		queryParams.Add("include", filters.Include)
	}

	var users []types.User = make([]types.User, 0)
	for {
		url.RawQuery = queryParams.Encode()
		req := &http.Request{
			Method: "GET",
			Header: headers,
			URL:    url,
		}
		req = req.WithContext(ctx)

		resp, err := ua.client.Do(req)
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

		var response types.ListUsersResponse
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}

		users = append(users, response.Users...)
		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		// increment the url params
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}

	return users, nil
}

func (ua *UserApi) CreateUser(ctx context.Context, request types.CreateUserRequest) (*types.User, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	headers := ua.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ua.baseUrl, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header = headers
	req = req.WithContext(ctx)

	resp, err := ua.client.Do(req)
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

	var user types.User
	if err := json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ua *UserApi) UpdateUser(ctx context.Context, request types.UpdateUserRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ua.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%d", ua.baseUrl, request.UserId), bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	req.Header = headers
	req = req.WithContext(ctx)

	resp, err := ua.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

func (ua *UserApi) DeleteUser(ctx context.Context, userId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", ua.baseUrl, userId), nil)
	if err != nil {
		return err
	}

	req.Header = ua.authHeader
	req = req.WithContext(ctx)

	resp, err := ua.client.Do(req)
	if err != nil {
		return err
	}

	if err := types.CheckResponse(resp); err != nil {
		return err
	}

	return nil
}
