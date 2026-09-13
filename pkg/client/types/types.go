package types

type PowerAction string

var (
	PowerActionStart   PowerAction = "start"
	PowerActionStop    PowerAction = "stop"
	PowerActionRestart PowerAction = "restart"
	PowerActionKill    PowerAction = "kill"
)

type ManagePowerRequest struct {
	ServerIdentifier string
	Signal           PowerAction `json:"signal"` // start, stop, restart, kill
}

type ListServersParams struct {
	Include string
	Page    int
	PerPage int
}

type ListServersResponse struct {
	Servers  []Server `json:"data"`
	MetaData `json:"meta"`
}

type GetServerActivityResponse struct {
	Data     []ActivityLog `json:"data"`
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

type ServerResources struct {
	Object     string `json:"object"`
	Attributes struct {
		CurrentState string `json:"current_state"`
		IsSuspended  bool   `json:"is_suspended"`
		Resources    struct {
			MemoryBytes    int64   `json:"memory_bytes"`
			CPUAbsolute    float64 `json:"cpu_absolute"`
			DiskBytes      int64   `json:"disk_bytes"`
			NetworkRxBytes int64   `json:"network_rx_bytes"`
			NetworkTxBytes int64   `json:"network_tx_bytes"`
			Uptime         int64   `json:"uptime"`
		} `json:"resources"`
	} `json:"attributes"`
}

type ActivityLog struct {
	Attributes struct {
		ID                    string         `json:"id"`
		Batch                 *string        `json:"batch"`
		Event                 string         `json:"event"`
		IsAPI                 bool           `json:"is_api"`
		IP                    string         `json:"ip"`
		Description           *string        `json:"description"`
		Properties            map[string]any `json:"properties"`
		HasAdditionalMetadata bool           `json:"has_additional_metadata"`
		Timestamp             string         `json:"timestamp"`
	} `json:"attributes"`
}

type ConsoleAccessDetails struct {
	Data struct {
		Token  string `json:"token"`
		Socket string `json:"socket"`
	} `json:"data"`
}
