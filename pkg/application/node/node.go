package node

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

type NodeApi struct {
	baseUrl    string
	client     *http.Client
	authHeader http.Header
}

func NewNodeApi(baseUrl string, client *http.Client, authHeader http.Header) *NodeApi {
	return &NodeApi{
		baseUrl:    baseUrl,
		client:     client,
		authHeader: authHeader,
	}
}

// NODES
func (ac *NodeApi) ListNodes(ctx context.Context, request *types.PteroListRequest) ([]types.Node, error) {
	headers := ac.authHeader
	url, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, "/api/application/nodes"))
	if err != nil {
		return nil, err
	}

	queryParams := request.BuildQueryParams(url)

	var nodes []types.Node
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

		var response ListNodesResponse

		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}
		nodes = append(nodes, response.Nodes...)

		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		// increment the url params
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}

	return nodes, nil
}

func (ac *NodeApi) GetNode(ctx context.Context, request GetNodeRequest) (*types.Node, error) {
	headers := ac.authHeader
	url, err := url.Parse(fmt.Sprintf("%s%s%s", ac.baseUrl, "/api/application/nodes/", strconv.Itoa(request.NodeId)))
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()

	if request.Include != nil {
		queryParams.Add("include", strings.Join(request.Include, ","))
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

	var node types.Node
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (ac *NodeApi) GetDeployableNodes(ctx context.Context, request GetDeployableNodesRequest) ([]types.Node, error) {
	headers := ac.authHeader
	url, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, "/api/application/nodes/deployable"))
	if err != nil {
		return nil, err
	}

	queryParams := url.Query()
	if request.Page != 0 {
		queryParams.Add("page", strconv.Itoa(request.Page))
	}
	if request.PerPage != 0 {
		queryParams.Add("per_page", strconv.Itoa(request.PerPage))
	}
	if request.Disk != 0 {
		queryParams.Add("disk", strconv.Itoa(request.Disk))
	}
	if request.Memory != 0 {
		queryParams.Add("memory", strconv.Itoa(request.Memory))
	}

	var nodes []types.Node
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

		var response ListNodesResponse

		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}
		nodes = append(nodes, response.Nodes...)

		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		// increment the url params
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}

	return nodes, nil
}

func (ac *NodeApi) CreateNode(ctx context.Context, request CreateNodeRequest) (*types.Node, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	headers := ac.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/application/nodes", ac.baseUrl), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header = headers
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

	var node types.Node
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (ac *NodeApi) UpdateNodeConfiguration(ctx context.Context, request UpdateNodeConfigurationRequest) error {
	if request.NodeId == 0 {
		return errors.New("missing required field 'NodeId'")
	}

	headers := ac.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/application/nodes/%d", ac.baseUrl, request.NodeId), bytes.NewBuffer(requestBytes))
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

	if resp.StatusCode >= 400 {
		var apiErr types.PterodactylAPIError
		if err := json.Unmarshal(b, &apiErr); err == nil {
			// Inspect returned errors
			for _, e := range apiErr.Errors {
				// If Wings failed to sync the config file, ignore and treat as success
				if e.Code == "ConfigurationNotPersistedException" {
					fmt.Println("Warning: Node updated in database, but Wings daemon config file could not be updated automatically.")
					return nil
				}
			}
		}

		return fmt.Errorf("error updating config - %s", string(b))
	}

	var node types.Node
	if err := json.Unmarshal(b, &node); err != nil {
		return err
	}

	return nil
}

func (ac *NodeApi) GetNodeConfiguration(ctx context.Context, nodeId int) (*types.NodeConfiguration, error) {
	headers := ac.authHeader
	url, err := url.Parse(fmt.Sprintf("%s/api/application/nodes/%s/configuration", ac.baseUrl, strconv.Itoa(nodeId)))
	if err != nil {
		return nil, err
	}

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

	var response types.NodeConfiguration
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type ListNodeAllocationsRequest struct {
	NodeId  int
	Page    int
	PerPage int
}

type ListNodeAllocationsResponse struct {
	Allocations    []types.Allocation `json:"data"`
	types.MetaData `json:"meta"`
}

func (ac *NodeApi) ListNodeAllocations(ctx context.Context, request ListNodeAllocationsRequest) ([]types.Allocation, error) {
	headers := ac.authHeader
	url, err := url.Parse(fmt.Sprintf("%s/api/application/nodes/%d/allocations", ac.baseUrl, request.NodeId))
	if err != nil {
		return nil, err
	}
	queryParams := url.Query()
	if request.Page != 0 {
		queryParams.Add("page", strconv.Itoa(request.Page))
	}
	if request.PerPage != 0 {
		queryParams.Add("per_page", strconv.Itoa(request.PerPage))
	}

	var allocations []types.Allocation = make([]types.Allocation, 0)
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

		var response ListNodeAllocationsResponse
		if err := json.Unmarshal(b, &response); err != nil {
			return nil, err
		}

		allocations = append(allocations, response.Allocations...)
		if response.MetaData.Paginaton.CurrentPage == response.MetaData.Paginaton.TotalPages {
			break
		}
		// increment the url params
		queryParams.Add("page", strconv.Itoa(response.MetaData.Paginaton.CurrentPage+1))
	}

	return allocations, nil
}

func (ac *NodeApi) CreateNodeAllocation(ctx context.Context, request CreateNodeAllocationRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.authHeader
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/application/nodes/%d/allocations", ac.baseUrl, request.NodeId), bytes.NewBuffer(requestBytes))
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error creating node allocation - %s", string(b))
	}

	return nil
}

func (ac *NodeApi) DeleteNodeAllocation(ctx context.Context, request DeleteNodeAllocationRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.authHeader
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/application/nodes/%d/allocations/%d", ac.baseUrl, request.NodeId, request.AllocationId), nil)
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error deleting node allocation - %s", string(b))
	}

	return nil
}

func (ac *NodeApi) DeleteNode(ctx context.Context, nodeId int) error {
	if nodeId <= 0 {
		return fmt.Errorf("missing required field '%s'", "NodeId")
	}

	headers := ac.authHeader
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/application/nodes/%d", ac.baseUrl, nodeId), nil)
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error deleting node - %s", string(b))
	}

	return nil
}
