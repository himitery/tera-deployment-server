package services

import (
	"errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strings"
	"tera/deployment/internal/domain/models"
	"tera/deployment/internal/ports"
	"tera/deployment/internal/usecases"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
)

type DeploymentManager struct {
	argocd   ports.Argocd
	services []config.ServiceConfig
}

func NewDeploymentManager(
	conf *config.Config,
	argocd ports.Argocd,
) usecases.DeploymentManager {
	return &DeploymentManager{
		argocd:   argocd,
		services: conf.Services,
	}
}

func (ctx *DeploymentManager) GetList() ([]models.ArgocdApplication, error) {
	applications, err := ctx.argocd.GetList()
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(applications, func(item models.ArgocdApplication, _ int) (models.ArgocdApplication, bool) {
		return item, ctx.hasService(item.Name)
	}), nil
}

func (ctx *DeploymentManager) Create(service, version, namespace string, values map[string]string) (*models.ArgocdApplication, error) {
	if namespace == "" {
		namespace = service
	}

	if !ctx.hasService(service) {
		logger.Warn("service not found", zap.String("service", service))

		return nil, errors.New("service not found")
	}

	application, err := ctx.argocd.Create(service, version, namespace, values)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (ctx *DeploymentManager) hasService(service string) bool {
	serviceNames := lo.Map(ctx.services, func(item config.ServiceConfig, _ int) string {
		return strings.ToLower(item.Name)
	})

	return lo.Contains(serviceNames, service)
}