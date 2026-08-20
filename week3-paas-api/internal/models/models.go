//this file contains json request and response models for the API.

package models

// CreateInstanceRequest is the JSON body used to create PostgreSQL.
type CreateInstanceRequest struct {
	Name        string `json:"name"`
	Instances   int64  `json:"instances"`
	StorageSize string `json:"storageSize"`
	Database    string `json:"database"`
	Owner       string `json:"owner"`
}

// Instance represents one PostgreSQL product instance.
type Instance struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// OwnedBy is the platform user who created the instance.
	// Owner below is the PostgreSQL role that owns the database.
	OwnedBy string `json:"ownedBy,omitempty"`

	Status           string `json:"status"`
	DesiredInstances int64  `json:"desiredInstances"`
	ReadyInstances   int64  `json:"readyInstances"`
	Primary          string `json:"primary,omitempty"`
	StorageSize      string `json:"storageSize"`
	Database         string `json:"database"`
	Owner            string `json:"owner"`
	CreatedAt        string `json:"createdAt,omitempty"`
}

// InstanceList is returned by GET /api/v1/instances.
type InstanceList struct {
	Items []Instance `json:"items"`
	Count int        `json:"count"`
}

// ConnectionData contains PostgreSQL connection credentials.
type ConnectionData struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	URI      string `json:"uri,omitempty"`
}

// OperationResponse describes an asynchronous operation.
type OperationResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ErrorResponse provides consistent REST API errors.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
