package types

import (
	"errors"
	"fmt"
)

type PteroListRequest struct {
	Include []string
	Page    int
	PerPage int
}

type ListUsersResponse struct {
	Users    []User `json:"data"`
	MetaData `json:"meta"`
}

type CreateNodeAllocationRequest struct {
	NodeId  int
	Ip      string   `json:"ip"`
	IpAlias string   `json:"ip_alias"`
	Ports   []string `json:"ports"`
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

type ListNodeAllocationsRequest struct {
	NodeId  int
	Page    int
	PerPage int
}

type ListNodeAllocationsResponse struct {
	Allocations []Allocation `json:"data"`
	MetaData    `json:"meta"`
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

type ListServersResponse struct {
	Servers  []Server `json:"data"`
	MetaData `json:"meta"`
}

type ListNodesResponse struct {
	Nodes    []Node `json:"data"`
	MetaData `json:"meta"`
}

type CreateServerAllocation struct {
	Default int `json:"default"`
	Backups int `json:"backups"`
}

type CreateServerDeploy struct {
	// required
	Locations   []int `json:"locations"`
	DedicatedIp bool  `json:"dedicated_ip"`
	// leave empty to dynamically grab empty available, required ports
	PortRange []int `json:"port_range"`
}

type CreateServerRequest struct {
	Name              string                  `json:"name"`
	User              int                     `json:"user"`
	Egg               int                     `json:"egg"`
	DockerImage       string                  `json:"docker_image,omitempty"`
	Startup           string                  `json:"startup,omitempty"`
	Environment       map[string]string       `json:"environment,omitempty"`
	Limits            Limits                  `json:"limits"`
	FeatureLimits     FeatureLimits           `json:"feature_limits"`
	Allocation        *CreateServerAllocation `json:"allocation"`
	Deploy            *CreateServerDeploy     `json:"deploy"`
	OOMDisabled       bool                    `json:"oom_disabled"`
	StartOnCompletion bool                    `json:"start_on_completion"`
}

func (r CreateServerRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.User == 0 {
		return errors.New("user is required")
	}
	if r.Egg == 0 {
		return errors.New("egg is required")
	}

	if r.Limits.Memory == 0 {
		return errors.New("limits.memory is required")
	}
	if r.Limits.Disk == 0 {
		return errors.New("limits.disk is required")
	}
	if r.Limits.IO == 0 {
		return errors.New("limits.io is required")
	}
	if r.Allocation != nil && r.Deploy != nil {
		return errors.New("allocation and deploy cannot both be specified, please choose one")
	}
	if r.Deploy != nil {
		if len(r.Deploy.Locations) <= 0 {
			return errors.New("deploy locations cannot be empty")
		}
	}

	return nil
}

type UpdateServerRequest interface {
	Validate() error
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

type CreateServerDatabaseRequest struct {
	ServerId               int
	Database               string `json:"database"`
	RemoteConnectionString string `json:"remote"` // can use % for all
	Host                   int    `json:"host"`
}

func (csdr CreateServerDatabaseRequest) Validate() error {
	if csdr.ServerId <= 0 {
		return errors.New("ServerId is required")
	}
	if csdr.Database == "" {
		return errors.New("database is required")
	}
	if csdr.RemoteConnectionString == "" {
		return errors.New("remote is required")
	}
	if csdr.Host <= 0 {
		return errors.New("host is required")
	}

	return nil
}

type UpdateServerDatabaseRequest struct {
	ServerId               int
	DatabaseId             int
	RemoteConnectionString string `json:"remote"` // can use % for all
}

func (usdr UpdateServerDatabaseRequest) Validate() error {
	if usdr.ServerId <= 0 {
		return errors.New("ServerId is required")
	}
	if usdr.DatabaseId <= 0 {
		return errors.New("DatabaseId is required")
	}

	return nil
}

type ListUsersFilters struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Email      string `json:"filter[email]"`
	UUID       string `json:"filter[uuid]"`
	Username   string `json:"filter[username]"`
	ExternalID string `json:"filter[external_id]"`
	Sort       string `json:"sort"`    // id, uuid, username, email, created_at, updated_at
	Include    string `json:"include"` // servers
}

type CreateUserRequest struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Password   string `json:"password"`
	Language   string `json:"language"`
	RootAdmin  bool   `json:"root_admin"`
	ExternalID string `json:"external_id"`
}

func (r *CreateUserRequest) Validate() error {
	// set defaults
	if r.Language == "" {
		r.Language = "en"
	}

	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.FirstName == "" {
		return errors.New("first_name is required")
	}
	if r.LastName == "" {
		return errors.New("last_name is required")
	}

	return nil
}

type UpdateUserRequest struct {
	UserId     int
	Email      string `json:"email"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Password   string `json:"password"`
	Language   string `json:"language"`
	RootAdmin  bool   `json:"root_admin"`
	ExternalID string `json:"external_id"`
}

func (r *UpdateUserRequest) Validate() error {
	// set defaults
	if r.Language == "" {
		r.Language = "en"
	}

	if r.UserId <= 0 {
		return errors.New("UserId is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.FirstName == "" {
		return errors.New("first_name is required")
	}
	if r.LastName == "" {
		return errors.New("last_name is required")
	}

	return nil
}

type ListNestEggsResponse struct {
	Eggs     []Egg `json:"data"`
	MetaData `json:"meta"`
}

type ListNestsResponse struct {
	Nests    []Nest `json:"data"`
	MetaData `json:"meta"`
}

type PterodactylAPIError struct {
	Errors []struct {
		Code   string `json:"code"`
		Status string `json:"status"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

type MetaData struct {
	Paginaton struct {
		Total       int `json:"total"`
		Count       int `json:"count"`
		PerPage     int `json:"per_page"`
		CurrentPage int `json:"current_page"`
		TotalPages  int `json:"total_pages"`
	} `json:"pagination"`
}
