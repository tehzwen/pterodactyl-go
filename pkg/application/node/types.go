package node

import (
	"fmt"

	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

type GetNodeRequest struct {
	NodeId  int
	Include []string
}

type GetDeployableNodesRequest struct {
	Page    int
	PerPage int
	Memory  int
	Disk    int
}

type CreateNodeRequest struct {
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	LocationId                   int    `json:"location_id"`
	FQDN                         string `json:"fqdn"`
	Scheme                       string `json:"scheme"`
	BehindProxy                  bool   `json:"behind_proxy"`
	Public                       bool   `json:"public"`
	DaemonBase                   string `json:"daemon_base"`
	DaemonSftp                   int    `json:"daemon_sftp"`
	DaemonListen                 int    `json:"daemon_listen"`
	Memory                       int    `json:"memory"`
	MemoryOverAllocatePercentage int    `json:"memory_overallocate"`
	Disk                         int    `json:"disk"`
	DiskOverAllocatePercentage   int    `json:"disk_overallocate"`
	MaxUploadSize                int    `json:"upload_size"`
	MaintenanceMode              bool   `json:"maintenance_mode"`
}

func (CNR *CreateNodeRequest) Validate() error {
	// check all required fields
	if CNR.Name == "" {
		return fmt.Errorf("missing required field '%s'", "name")
	}
	if CNR.LocationId == 0 {
		return fmt.Errorf("missing required field '%s'", "location_id")
	}
	if CNR.FQDN == "" {
		return fmt.Errorf("missing required field '%s'", "fqdn")
	}
	if CNR.Memory <= 0 {
		return fmt.Errorf("missing required field '%s'", "memory")
	}
	if CNR.Disk <= 0 {
		return fmt.Errorf("missing required field '%s'", "disk")
	}
	if CNR.Scheme == "" {
		return fmt.Errorf("missing required field '%s'", "scheme")
	}
	// set defaults if they aren't already set
	if CNR.DaemonSftp == 0 {
		CNR.DaemonSftp = 2022
	}
	if CNR.DaemonListen == 0 {
		CNR.DaemonListen = 8080
	}
	if CNR.MaxUploadSize == 0 {
		CNR.MaxUploadSize = 100
	}
	if CNR.DaemonBase == "" {
		CNR.DaemonBase = "/var/lib/pterodactyl/volumes"
	}

	return nil
}

type UpdateNodeConfigurationRequest struct {
	NodeId                       int
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	LocationId                   int    `json:"location_id"`
	FQDN                         string `json:"fqdn"`
	Scheme                       string `json:"scheme"`
	BehindProxy                  bool   `json:"behind_proxy"`
	Public                       bool   `json:"public"`
	DaemonBase                   string `json:"daemon_base"`
	DaemonSftp                   int    `json:"daemon_sftp"`
	DaemonListen                 int    `json:"daemon_listen"`
	Memory                       int    `json:"memory"`
	MemoryOverAllocatePercentage int    `json:"memory_overallocate"`
	Disk                         int    `json:"disk"`
	DiskOverAllocatePercentage   int    `json:"disk_overallocate"`
	MaxUploadSize                int    `json:"upload_size"`
	MaintenanceMode              bool   `json:"maintenance_mode"`
}

type DeleteNodeAllocationRequest struct {
	NodeId       int
	AllocationId int
}

func (dar *DeleteNodeAllocationRequest) Validate() error {
	if dar.NodeId <= 0 {
		return fmt.Errorf("missing required field '%s'", "NodeId")
	}
	if dar.AllocationId <= 0 {
		return fmt.Errorf("missing required field '%s'", "AllocationId")
	}
	return nil
}

type CreateNodeAllocationRequest struct {
	NodeId  int
	Ip      string   `json:"ip"`
	IpAlias string   `json:"ip_alias"`
	Ports   []string `json:"ports"`
}

type ListNodesResponse struct {
	Nodes          []types.Node `json:"data"`
	types.MetaData `json:"meta"`
}

func (car *CreateNodeAllocationRequest) Validate() error {
	if car.Ip == "" {
		return fmt.Errorf("missing required field '%s'", "ip")
	}
	if len(car.Ports) == 0 {
		return fmt.Errorf("missing required field '%s'", "ports")
	}
	if car.NodeId == 0 {
		return fmt.Errorf("missing required field '%s'", "NodeId")
	}

	return nil
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
