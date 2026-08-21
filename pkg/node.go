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

// NODES
func (ac *ApplicationClient) ListNodes(ctx context.Context, request *PteroListRequest) ([]Node, error) {
	headers := ac.buildHeaders()
	url, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, "/api/application/nodes"))
	if err != nil {
		return nil, err
	}

	queryParams := request.buildQueryParams(url)

	var nodes []Node
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

func (ac *ApplicationClient) GetNode(ctx context.Context, request GetNodeRequest) (*Node, error) {
	headers := ac.buildHeaders()
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

	var node Node
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (ac *ApplicationClient) GetDeployableNodes(ctx context.Context, request GetDeployableNodesRequest) ([]Node, error) {
	headers := ac.buildHeaders()
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

	var nodes []Node
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

func (ac *ApplicationClient) CreateNode(ctx context.Context, request CreateNodeRequest) (*Node, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	headers := ac.buildHeaders()
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

	var node Node
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (ac *ApplicationClient) UpdateNodeConfiguration(ctx context.Context, request UpdateNodeConfigurationRequest) error {
	if request.NodeId == 0 {
		return errors.New("missing required field 'NodeId'")
	}

	headers := ac.buildHeaders()
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
		var apiErr PterodactylAPIError
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

	var node Node
	if err := json.Unmarshal(b, &node); err != nil {
		return err
	}

	return nil
}

func (ac *ApplicationClient) GetNodeConfiguration(ctx context.Context, nodeId int) (*NodeConfiguration, error) {
	headers := ac.buildHeaders()
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

	var response NodeConfiguration
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
	Allocations []Allocation `json:"data"`
	MetaData    `json:"meta"`
}

func (ac *ApplicationClient) ListNodeAllocations(ctx context.Context, request ListNodeAllocationsRequest) ([]Allocation, error) {
	headers := ac.buildHeaders()
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

	var allocations []Allocation = make([]Allocation, 0)
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

func (ac *ApplicationClient) CreateNodeAllocation(ctx context.Context, request CreateNodeAllocationRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.buildHeaders()
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

func (ac *ApplicationClient) DeleteNodeAllocation(ctx context.Context, request DeleteNodeAllocationRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers := ac.buildHeaders()
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

func (ac *ApplicationClient) DeleteNode(ctx context.Context, nodeId int) error {
	if nodeId <= 0 {
		return fmt.Errorf("missing required field '%s'", "NodeId")
	}

	headers := ac.buildHeaders()
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
