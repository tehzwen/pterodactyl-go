package types

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PteroListRequest struct {
	Include []string
	Page    int
	PerPage int
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

type CreateNodeAllocationRequest struct {
	NodeId  int
	Ip      string   `json:"ip"`
	IpAlias string   `json:"ip_alias"`
	Ports   []string `json:"ports"`
}

type ListUsersResponse struct {
	Users    []User `json:"data"`
	MetaData `json:"meta"`
}

type ListServersResponse struct {
	Servers  []Server `json:"data"`
	MetaData `json:"meta"`
}

type ListNodesResponse struct {
	Nodes    []Node `json:"data"`
	MetaData `json:"meta"`
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

type CreateServerRequest struct {
	Name          string            `json:"name"`
	User          int               `json:"user"`
	Egg           int               `json:"egg"`
	DockerImage   string            `json:"docker_image,omitempty"`
	Startup       string            `json:"startup,omitempty"`
	Environment   map[string]string `json:"environment,omitempty"`
	Limits        Limits            `json:"limits"`
	FeatureLimits FeatureLimits     `json:"feature_limits"`
	Allocation    struct {
		Default int `json:"default"`
		Backups int `json:"backups"`
	} `json:"allocation"`
	Deploy map[string]string `json:"deploy,omitempty"`
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

type User struct {
	Attributes struct {
		ID         int       `json:"id"`
		ExternalID *string   `json:"external_id"`
		UUID       string    `json:"uuid"`
		Username   string    `json:"username"`
		Email      string    `json:"email"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		Language   string    `json:"language"`
		RootAdmin  bool      `json:"root_admin"`
		TwoFactor  bool      `json:"2fa"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"attributes"`
}

type Node struct {
	Attributes struct {
		ID                 int       `json:"id"`
		UUID               string    `json:"uuid"`
		Public             bool      `json:"public"`
		Name               string    `json:"name"`
		Description        string    `json:"description"`
		LocationID         int       `json:"location_id"`
		FQDN               string    `json:"fqdn"`
		Scheme             string    `json:"scheme"`
		BehindProxy        bool      `json:"behind_proxy"`
		MaintenanceMode    bool      `json:"maintenance_mode"`
		Memory             int       `json:"memory"`
		MemoryOverallocate int       `json:"memory_overallocate"`
		Disk               int       `json:"disk"`
		DiskOverallocate   int       `json:"disk_overallocate"`
		UploadSize         int       `json:"upload_size"`
		DaemonListen       int       `json:"daemon_listen"`
		DaemonSFTP         int       `json:"daemon_sftp"`
		DaemonBase         string    `json:"daemon_base"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		AllocatedResources struct {
			Memory int `json:"memory"`
			Disk   int `json:"disk"`
		} `json:"allocated_resources"`
	} `json:"attributes"`
}

func (n *Node) GetUpdateRequest() UpdateNodeConfigurationRequest {
	return UpdateNodeConfigurationRequest{
		NodeId:                       n.Attributes.ID,
		Name:                         n.Attributes.Name,
		Description:                  n.Attributes.Description,
		LocationId:                   n.Attributes.LocationID,
		FQDN:                         n.Attributes.FQDN,
		Scheme:                       n.Attributes.Scheme,
		BehindProxy:                  n.Attributes.BehindProxy,
		Public:                       n.Attributes.Public,
		DaemonBase:                   n.Attributes.DaemonBase,
		DaemonSftp:                   n.Attributes.DaemonSFTP,
		DaemonListen:                 n.Attributes.DaemonListen,
		Memory:                       n.Attributes.Memory,
		MemoryOverAllocatePercentage: n.Attributes.MemoryOverallocate,
		Disk:                         n.Attributes.Disk,
		DiskOverAllocatePercentage:   n.Attributes.DiskOverallocate,
		MaxUploadSize:                n.Attributes.UploadSize,
		MaintenanceMode:              n.Attributes.MaintenanceMode,
	}
}

type NodeConfiguration struct {
	Debug         bool     `json:"debug"`
	UUID          string   `json:"uuid"`
	TokenID       string   `json:"token_id"`
	Token         string   `json:"token"`
	AllowedMounts []string `json:"allowed_mounts"`
	Remote        string   `json:"remote"`

	API struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		UploadLimit int64  `json:"upload_limit"`
		SSL         struct {
			Enabled bool   `json:"enabled"`
			Cert    string `json:"cert"`
			Key     string `json:"key"`
		} `json:"ssl"`
	} `json:"api"`

	System struct {
		Data string `json:"data"`
		SFTP struct {
			BindPort int `json:"bind_port"`
		} `json:"sftp"`
	} `json:"system"`
}

type Server struct {
	Attributes struct {
		Id          int    `json:"id"`
		ExternalId  string `json:"external_id"`
		UUID        string `json:"uuid"`
		Identifier  string `json:"identifier"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Suspended   bool   `json:"suspended"`
		Limits      struct {
			MemoryMB    int  `json:"number"`
			SwapMB      int  `json:"swap"`
			DiskMB      int  `json:"disk"`
			Io          int  `json:"io"`
			CPU         int  `json:"cpu"`
			Threads     int  `json:"threads"` // confirm this field, null at the time
			OOMDisabled bool `json:"oom_disabled"`
		} `json:"limits"`
		FeatureLimits struct {
			Databases   int `json:"databases"`
			Allocations int `json:"allocations"`
			Backups     int `json:"backups"`
		} `json:"feature_limits"`
		User       int `json:"user"`
		Node       int `json:"node"`
		Allocation int `json:"allocation"`
		Nest       int `json:"nest"`
		Egg        int `json:"egg"`
		Container  struct {
			StartupCommand string         `json:"startup_command"`
			Image          string         `json:"image"`
			Installed      int            `json:"installed"`
			Environment    map[string]any `json:"environment"`
			SkipScripts    bool           `json:"skip_scripts"`
		} `json:"container"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		// relationships
		Relationships struct {
			Node *Node `json:"node"`
			User *User `json:"user"`
		} `json:"relationships"`
	} `json:"attributes"`
}

type Allocation struct {
	Attributes struct {
		ID       int     `json:"id"`
		IP       string  `json:"ip"`
		IPAlias  *string `json:"ip_alias"`
		Port     int     `json:"port"`
		Notes    *string `json:"notes"`
		Assigned bool    `json:"assigned"`
	} `json:"attributes"`
}

type Limits struct {
	Memory      int    `json:"memory"`
	Swap        int    `json:"swap"`
	Disk        int    `json:"disk"`
	IO          int    `json:"io"`
	CPU         int    `json:"cpu"`
	Threads     string `json:"threads,omitempty"`
	OOMDisabled bool   `json:"oom_disabled,omitempty"`
}

type FeatureLimits struct {
	Databases   int `json:"databases"`
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
}

func (pt *PteroListRequest) BuildQueryParams(url *url.URL) url.Values {
	queryParams := url.Query()

	if pt.Include != nil {
		queryParams.Add("include", strings.Join(pt.Include, ","))
	}

	if pt.Page != 0 {
		queryParams.Add("page", strconv.Itoa(pt.Page))
	}

	if pt.PerPage != 0 {
		queryParams.Add("per_page", strconv.Itoa(pt.PerPage))
	}
	return queryParams
}

type Database struct {
	ID             int       `json:"id"`
	Server         int       `json:"server"`
	Host           int       `json:"host"`
	Database       string    `json:"database"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Remote         string    `json:"remote"`
	MaxConnections int       `json:"max_connections"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Relationships  struct {
		Host struct {
			Object     string `json:"object"`
			Attributes struct {
				ID           int       `json:"id"`
				Name         string    `json:"name"`
				Host         string    `json:"host"`
				Port         int       `json:"port"`
				Username     string    `json:"username"`
				MaxDatabases int       `json:"max_databases"`
				CreatedAt    time.Time `json:"created_at"`
				UpdatedAt    time.Time `json:"updated_at"`
			} `json:"attributes"`
		} `json:"host"`
	} `json:"relationships"`
}
