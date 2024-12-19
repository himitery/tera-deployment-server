package argocd

import (
	"context"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/samber/lo"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"tera/deployment/internal/domain/models"
	"tera/deployment/internal/ports"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
)

type Argocd struct {
	client        apiclient.Client
	repository    string
	metaNamespace string
}

func NewArgocd(conf *config.Config) ports.Argocd {
	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: conf.Argocd.URL,
		AuthToken:  conf.Argocd.Token,
		GRPCWeb:    true,
	})
	if err != nil {
		logger.Error("failed to create Argocd client")

		panic(err)
	}

	return &Argocd{
		client:        client,
		repository:    conf.Argocd.Repository,
		metaNamespace: conf.Argocd.Metadata.Namespace,
	}
}

func (ctx *Argocd) GetList() ([]models.ArgocdApplication, error) {
	io, client, err := ctx.client.NewApplicationClient()
	if err != nil {
		logger.Error("failed to create Argocd application client", zap.Error(err))

		return nil, err
	}
	defer io.Close()

	data, err := client.List(context.Background(), &application.ApplicationQuery{})
	if err != nil {
		logger.Error("failed to list Argocd applications", zap.Error(err))

		return nil, err
	}

	return lo.Map(data.Items, func(item v1alpha1.Application, index int) models.ArgocdApplication {
		return models.ArgocdApplication{
			Name:    strings.ToLower(item.Name),
			Version: item.Spec.Source.TargetRevision,
		}
	}), nil
}

func (ctx *Argocd) Create(service, version, namespace string, values map[string]string) (*models.ArgocdApplication, error) {
	io, client, err := ctx.client.NewApplicationClient()
	if err != nil {
		logger.Error("failed to create Argocd application client", zap.Error(err))

		return nil, err
	}
	defer io.Close()

	parameters := lo.MapToSlice(values, func(key string, value string) v1alpha1.HelmParameter {
		return v1alpha1.HelmParameter{
			Name:        key,
			Value:       value,
			ForceString: false,
		}
	})

	data, err := client.Create(context.Background(), &application.ApplicationCreateRequest{
		Application: &v1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      service,
				Namespace: ctx.metaNamespace,
			},
			Spec: v1alpha1.ApplicationSpec{
				Project: "default",
				Source: &v1alpha1.ApplicationSource{
					RepoURL:        ctx.repository,
					Chart:          service,
					TargetRevision: version,
					Helm: &v1alpha1.ApplicationSourceHelm{
						ReleaseName: service,
						Namespace:   namespace,
						Parameters:  parameters,
					},
				},
				Destination: v1alpha1.ApplicationDestination{
					Server:    "https://kubernetes.default.svc",
					Namespace: namespace,
				},
				SyncPolicy: &v1alpha1.SyncPolicy{
					Automated: &v1alpha1.SyncPolicyAutomated{
						Prune:      true,
						SelfHeal:   true,
						AllowEmpty: false,
					},
					Retry: &v1alpha1.RetryStrategy{
						Backoff: &v1alpha1.Backoff{
							Duration:    "5s",
							Factor:      lo.ToPtr(int64(2)),
							MaxDuration: "3m",
						},
						Limit: 5,
					},
					SyncOptions: []string{
						"CreateNamespace=true",
						"ApplyOutOfSyncOnly=true",
						"ServerSideApply=true",
					},
				},
			},
		},
		Upsert:   lo.ToPtr(false),
		Validate: lo.ToPtr(true),
	})
	if err != nil {
		logger.Error("failed to create Argocd application", zap.Error(err))

		return nil, err
	}

	err = ctx.waitForApplicationSync(service, time.Minute*3)
	if err != nil {
		logger.Error("Argocd.Create: application sync failed", zap.Error(err))
		return nil, err
	}

	return &models.ArgocdApplication{
		Name:    strings.ToLower(data.Name),
		Version: data.Spec.Source.TargetRevision,
	}, nil
}

func (ctx *Argocd) waitForApplicationSync(service string, timeout time.Duration) error {
	syncCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	io, client, err := ctx.client.NewApplicationClient()
	if err != nil {
		logger.Error("failed to create Argocd application client", zap.Error(err))

		return err
	}
	defer io.Close()

	for {
		select {
		case <-syncCtx.Done():
			return syncCtx.Err()
		case <-ticker.C:
			data, err := client.Get(context.Background(), &application.ApplicationQuery{
				Name: &service,
			})
			if err != nil {
				logger.Error("failed to get Argocd application", zap.Error(err))

				continue
			}

			logger.Info(
				"Argocd application status",
				zap.Any("status", data.Status.Sync.Status),
				zap.Any("healthStatus", data.Status.Health.Status),
			)

			if data.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced && data.Status.Health.Status == health.HealthStatusHealthy {
				logger.Info("Argocd application synced")

				return nil
			}
		}
	}
}
