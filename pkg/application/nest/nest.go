package nest

import (
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

type NestApi struct {
	baseUrl    string
	client     *http.Client
	authHeader http.Header
}

func NewNestApi(baseUrl string, client *http.Client, authHeader http.Header) *NestApi {
	return &NestApi{
		baseUrl:    baseUrl,
		client:     client,
		authHeader: authHeader.Clone(),
	}
}

// ListNests lists existing nests defined in the system. The request can be set to `eggs,servers`
// to include these in the response
func (na *NestApi) ListNests(ctx context.Context, request *types.PteroListRequest) ([]types.Nest, error) {
	headers := na.authHeader.Clone()
	url, err := url.Parse(na.baseUrl)
	if err != nil {
		return nil, err
	}

	queryParams := request.BuildQueryParams(url)
	var nests []types.Nest = make([]types.Nest, 0)

	for {
		url.RawQuery = queryParams.Encode()
		req := &http.Request{
			Method: "GET",
			Header: headers,
			URL:    url,
		}
		req = req.WithContext(ctx)

		resp, err := na.client.Do(req)
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

		var response types.ListNestsResponse
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}

		nests = append(nests, response.Nests...)
		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		// increment the url params
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}
	return nests, nil
}

func (na *NestApi) GetNestDetails(ctx context.Context, nestId int, include []string) (*types.Nest, error) {
	headers := na.authHeader.Clone()
	baseUrl := fmt.Sprintf("%s/%d", na.baseUrl, nestId)
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
	resp, err := na.client.Do(req)
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

	var nest types.Nest
	if err := json.Unmarshal(b, &nest); err != nil {
		return nil, err
	}

	return &nest, nil
}

func (na *NestApi) ListNestEggs(ctx context.Context, nestId int, request *types.PteroListRequest) ([]types.Egg, error) {
	headers := na.authHeader.Clone()
	url, err := url.Parse(fmt.Sprintf("%s/%d/eggs", na.baseUrl, nestId))
	if err != nil {
		return nil, err
	}

	queryParams := request.BuildQueryParams(url)
	var eggs []types.Egg = make([]types.Egg, 0)

	for {
		url.RawQuery = queryParams.Encode()
		req := &http.Request{
			Method: "GET",
			Header: headers,
			URL:    url,
		}
		req = req.WithContext(ctx)

		resp, err := na.client.Do(req)
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

		var response types.ListNestEggsResponse
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}

		eggs = append(eggs, response.Eggs...)
		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}
	return eggs, nil
}

func (na *NestApi) GetEggDetails(ctx context.Context, nestId, eggId int, include []string) (*types.Egg, error) {
	headers := na.authHeader.Clone()
	baseUrl := fmt.Sprintf("%s/%d/eggs/%d", na.baseUrl, nestId, eggId)
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
	resp, err := na.client.Do(req)
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

	var egg types.Egg
	if err := json.Unmarshal(b, &egg); err != nil {
		return nil, err
	}

	return &egg, nil
}
