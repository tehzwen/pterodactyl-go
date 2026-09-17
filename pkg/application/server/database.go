package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

func (sa *ServerApi) ListDatabases(ctx context.Context, serverId int) ([]types.Database, error) {
	headers := sa.authHeader.Clone()
	url, err := url.Parse(fmt.Sprintf("%s/%d/databases", sa.baseUrl, serverId))
	if err != nil {
		return nil, err
	}

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

	var databases []types.Database
	if err := json.Unmarshal(b, &databases); err != nil {
		return nil, err
	}

	return databases, nil
}

func (sa *ServerApi) CreateDatabase(ctx context.Context, request types.CreateServerDatabaseRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := sa.authHeader.Clone()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%d/databases", sa.baseUrl, request.ServerId), bytes.NewBuffer(requestBytes))
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

func (sa *ServerApi) UpdateDatabase(ctx context.Context, request types.UpdateServerDatabaseRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := sa.authHeader.Clone()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%d/databases/%d", sa.baseUrl, request.ServerId, request.DatabaseId), bytes.NewBuffer(requestBytes))
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

func (sa *ServerApi) ResetDatabasePassword(ctx context.Context, serverId, databaseId int) error {
	headers := sa.authHeader.Clone()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%d/databases/%d/reset-password", sa.baseUrl, serverId, databaseId), nil)
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

func (sa *ServerApi) DeleteDatabase(ctx context.Context, serverId, databaseId int) error {
	headers := sa.authHeader.Clone()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d/databases/%d", sa.baseUrl, serverId, databaseId), nil)
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
