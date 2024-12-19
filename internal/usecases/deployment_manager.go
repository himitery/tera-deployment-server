package usecases

import "tera/deployment/internal/domain/models"

type DeploymentManager interface {
	GetList() ([]models.Application, error)
	Create(service, version, namespace string, values map[string]string) (*models.Application, error)
}
