package models

type EventMessage struct {
	Action    string            `json:"action"` // create, delete
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Namespace string            `json:"namespace"`
	Values    map[string]string `json:"values"`
}
