package services

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/storage"
	"akshidas/e-com/pkg/types"
)

type ResourceModeler interface {
	GetAll() ([]*types.Resource, error)
	GetOne(int) (*types.Resource, error)
	Create(*types.CreateResourceRequest) error
	Update(int, *types.CreateResourceRequest) (*types.Resource, error)
	Delete(int) error
}

type ResourceService struct {
	roleModel ResourceModeler
}

func (r *ResourceService) GetAll() ([]*types.Resource, error) {
	return r.roleModel.GetAll()
}

func (r *ResourceService) GetOne(id int) (*types.Resource, error) {
	return r.roleModel.GetOne(id)
}

func (r *ResourceService) Create(newResource *types.CreateResourceRequest) error {
	return r.roleModel.Create(newResource)
}

func (r *ResourceService) Update(id int, newResource *types.CreateResourceRequest) (*types.Resource, error) {
	return r.roleModel.Update(id, newResource)
}

func (r *ResourceService) Delete(id int) error {
	return r.roleModel.Delete(id)
}

func NewResourceService(database *db.Storage) *ResourceService {
	roleModel := storage.NewResourceStorage(database.DB)
	return &ResourceService{
		roleModel: roleModel,
	}
}
