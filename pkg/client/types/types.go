package types

type ListServersParams struct {
	Include string
	Page    int
	PerPage int
}

type ListServersResponse struct {
	Servers  []Server `json:"data"`
	MetaData `json:"meta"`
}

type Server struct {
	Attributes struct {
		ServerOwner     bool                `json:"server_owner"`
		Identifier      string              `json:"identifier"`
		InternalID      int64               `json:"internal_id"`
		UUID            string              `json:"uuid"`
		Name            string              `json:"name"`
		Description     string              `json:"description"`
		Status          *string             `json:"status"`
		IsSuspended     bool                `json:"is_suspended"`
		IsInstalling    bool                `json:"is_installing"`
		IsTransferring  bool                `json:"is_transferring"`
		Node            string              `json:"node"`
		SFTPDetails     SFTPDetails         `json:"sftp_details"`
		Invocation      string              `json:"invocation"`
		DockerImage     string              `json:"docker_image"`
		EggFeatures     []string            `json:"egg_features"`
		FeatureLimits   FeatureLimits       `json:"feature_limits"`
		UserPermissions []string            `json:"user_permissions"`
		Limits          ServerLimits        `json:"limits"`
		Relationships   ServerRelationships `json:"relationships"`
	} `json:"attributes"`
}

type SFTPDetails struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type FeatureLimits struct {
	Databases   int `json:"databases"`
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
}

type ServerLimits struct {
	Memory  int64   `json:"memory"`
	Swap    int64   `json:"swap"`
	Disk    int64   `json:"disk"`
	IO      int64   `json:"io"`
	CPU     int64   `json:"cpu"`
	Threads *string `json:"threads"`
}

type ServerRelationships struct {
	Allocations AllocationList `json:"allocations"`
	Variables   VariableList   `json:"variables"`
}

type AllocationList struct {
	Object string              `json:"object"`
	Data   []AllocationElement `json:"data"`
}

type AllocationElement struct {
	Object     string               `json:"object"`
	Attributes AllocationAttributes `json:"attributes"`
}

type AllocationAttributes struct {
	ID        int64   `json:"id"`
	IP        string  `json:"ip"`
	IPAlias   *string `json:"ip_alias"`
	Port      int     `json:"port"`
	Notes     *string `json:"notes"`
	IsDefault bool    `json:"is_default"`
}

type VariableList struct {
	Object string            `json:"object"`
	Data   []VariableElement `json:"data"`
}

type VariableElement struct {
	Object     string             `json:"object"`
	Attributes VariableAttributes `json:"attributes"`
}

type VariableAttributes struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	EnvVariable  string `json:"env_variable"`
	DefaultValue string `json:"default_value"`
	ServerValue  string `json:"server_value"`
	IsEditable   bool   `json:"is_editable"`
	Rules        string `json:"rules"`
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
