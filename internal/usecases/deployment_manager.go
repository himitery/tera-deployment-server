package usecases

import "tera/deployment/internal/domain/models"

type DeploymentManager interface {
	GetList() ([]models.ArgocdApplication, error)
	Create(service, version, namespace string, values map[string]string) (*models.ArgocdApplication, error)
}
