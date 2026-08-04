package main

type CreateInstanceRequest struct {
	Name        string `json:"name"`
	Instances   int64  `json:"instances,omitempty"`
	StorageSize string `json:"storageSize,omitempty"`
	Database    string `json:"database,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

type Instance struct {
	Name             string `json:"name"`
	Namespace        string `json:"namespace"`
	Status           string `json:"status"`
	StatusMessage    string `json:"statusMessage,omitempty"`
	DesiredInstances int64  `json:"desiredInstances"`
	ReadyInstances   int64  `json:"readyInstances"`
	Primary          string `json:"primary,omitempty"`
	StorageSize      string `json:"storageSize"`
	Database         string `json:"database"`
	Owner            string `json:"owner"`
	CreatedAt        string `json:"createdAt,omitempty"`
}

type InstanceList struct {
	Items []Instance `json:"items"`
	Count int        `json:"count"`
}

type ConnectionData struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	URI      string `json:"uri,omitempty"`
}

type OperationResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
