package server

import (
	"errors"

	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
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
