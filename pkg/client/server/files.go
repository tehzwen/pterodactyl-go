package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tehzwen/pterodactyl-go/pkg/client/types"
)

func (sa *ServerApi) ListServerFiles(ctx context.Context, params types.ListServerFilesParams) (*types.ListServerFilesResponse, error) {
	if params.ServerId == "" {
		return nil, errors.New("list files request must include ServerId as param")
	}

	headers := sa.authHeader.Clone()
	url, err := url.Parse(fmt.Sprintf("%s/servers/%s/files/list", sa.baseUrl, params.ServerId))
	if err != nil {
		return nil, err
	}

	queryParams := url.Query()
	if params.FolderPath != "" {
		queryParams.Add("directory", params.FolderPath)
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

	var response types.ListServerFilesResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (sa *ServerApi) GetServerFileContents(ctx context.Context, params types.GetServerFileContentsParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}

	headers := sa.authHeader.Clone()
	url, err := url.Parse(fmt.Sprintf("%s/servers/%s/files/contents", sa.baseUrl, params.ServerId))
	if err != nil {
		return "", err
	}

	queryParams := url.Query()
	if params.File != "" {
		queryParams.Add("file", params.File)
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
		return "", err
	}

	if err := types.CheckResponse(resp); err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(b), nil
}
