package models

type Key struct {
	Value string
}

var (
	ArgocdApplicationList   Key = Key{Value: "ArgocdApplicationList"}
	ArgocdApplicationStatus Key = Key{Value: "ArgocdApplicationStatus"}
)

var (
	FetchArgocdApplication      = Key{Value: "FetchArgocdApplication"}
	CreateArgocdApplication Key = Key{Value: "CreateArgocdApplication"}
)
