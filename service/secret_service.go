package service

import (
	"context"

	"github.com/odpf/optimus/models"
	"github.com/odpf/optimus/store"
)

type SecretService interface {
	Save(context.Context, string, string, models.ProjectSecretItem) error
	Update(context.Context, string, string, models.ProjectSecretItem) error
	List(context.Context, string) ([]models.SecretItemInfo, error)
	Delete(context.Context, string, string, string) error
}

type SecretRepoFactory interface {
	New(projectSpec models.ProjectSpec) store.ProjectSecretRepository
}

type secretService struct {
	projService   ProjectService
	nsService     NamespaceService
	secretRepoFac SecretRepoFactory
}

func NewSecretService(projectService ProjectService, namespaceService NamespaceService, factory SecretRepoFactory) *secretService {
	return &secretService{
		projService:   projectService,
		nsService:     namespaceService,
		secretRepoFac: factory,
	}
}

func (s secretService) Save(ctx context.Context, projectName string, namespaceName string, item models.ProjectSecretItem) error {
	if item.Name == "" {
		return NewError(models.SecretEntity, ErrInvalidArgument, "secret name cannot be empty")
	}

	proj, namespace, err := s.nsService.GetNamespaceOptionally(ctx, projectName, namespaceName)
	if err != nil {
		return err
	}

	repo := s.secretRepoFac.New(proj)
	err = repo.Save(ctx, namespace, item)
	if err != nil {
		return FromError(err, models.SecretEntity, "error while saving secret")
	}
	return nil
}

func (s secretService) Update(ctx context.Context, projectName string, namespaceName string, item models.ProjectSecretItem) error {
	if item.Name == "" {
		return NewError(models.SecretEntity, ErrInvalidArgument, "secret name cannot be empty")
	}

	proj, namespace, err := s.nsService.GetNamespaceOptionally(ctx, projectName, namespaceName)
	if err != nil {
		return err
	}

	repo := s.secretRepoFac.New(proj)
	err = repo.Update(ctx, namespace, item)
	if err != nil {
		return FromError(err, models.SecretEntity, "error while updating secret")
	}
	return nil
}

func (s secretService) List(ctx context.Context, projectName string) ([]models.SecretItemInfo, error) {
	projectSpec, err := s.projService.Get(ctx, projectName)
	if err != nil {
		return nil, err
	}

	repo := s.secretRepoFac.New(projectSpec)
	secretItems, err := repo.GetAll(ctx)
	if err != nil {
		return []models.SecretItemInfo{}, FromError(err, models.SecretEntity, "error while saving secret")
	}
	return secretItems, nil
}

func (s secretService) Delete(ctx context.Context, projectName, namespaceName, secretName string) error {
	proj, namespace, err := s.nsService.GetNamespaceOptionally(ctx, projectName, namespaceName)
	if err != nil {
		return err
	}

	repo := s.secretRepoFac.New(proj)
	err = repo.Delete(ctx, namespace, secretName)
	if err != nil {
		return FromError(err, models.SecretEntity, "error while deleting secret")
	}
	return nil
}
