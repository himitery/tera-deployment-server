package ports

import "tera/deployment/internal/domain/models"

type Argocd interface {
	GetList() ([]models.Application, error)

	Create(service, version, namespace string, values map[string]string) (*models.Application, error)
}
