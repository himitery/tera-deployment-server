package models

type Key struct {
	Value string
}

var (
	ArgocdApplicationList   Key = Key{Value: "argocd_application_list"}
	ArgocdApplicationStatus Key = Key{Value: "argocd_application_status"}
)

var (
	FetchArgocdApplication      = Key{Value: "fetch_argocd_application"}
	CreateArgocdApplication Key = Key{Value: "create_argocd_application"}
)
