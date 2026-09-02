package server

import (
	"errors"

	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

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
	Servers        []types.Server `json:"data"`
	types.MetaData `json:"meta"`
}

type ListNodesResponse struct {
	Nodes          []types.Node `json:"data"`
	types.MetaData `json:"meta"`
}

type CreateServerRequest struct {
	Name          string              `json:"name"`
	User          int                 `json:"user"`
	Egg           int                 `json:"egg"`
	DockerImage   string              `json:"docker_image,omitempty"`
	Startup       string              `json:"startup,omitempty"`
	Environment   map[string]string   `json:"environment,omitempty"`
	Limits        types.Limits        `json:"limits"`
	FeatureLimits types.FeatureLimits `json:"feature_limits"`
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
