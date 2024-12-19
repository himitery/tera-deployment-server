package services

import (
	"errors"
	"fmt"
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

	if depends := ctx.findDepends(service); len(depends) > 0 {
		message := fmt.Sprintf(
			"service '%s' cannot be installed because the following dependencies are missing: %v",
			service,
			depends,
		)

		logger.Warn(message)

		return nil, errors.New(message)
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

func (ctx *DeploymentManager) findDepends(service string) []models.ArgocdApplication {
	application, _ := lo.Find(ctx.services, func(item config.ServiceConfig) bool {
		return item.Name == service
	})

	deployedApplications, _ := ctx.argocd.GetList()
	deployedApplicationNames := lo.Map(deployedApplications, func(item models.ArgocdApplication, _ int) string {
		return item.Name
	})

	return lo.FilterMap(application.Depends, func(item config.ServiceDependConfig, _ int) (models.ArgocdApplication, bool) {
		return models.ArgocdApplication{
			Name:    item.Name,
			Version: item.Version,
		}, !lo.Contains(deployedApplicationNames, item.Name)
	})
}
