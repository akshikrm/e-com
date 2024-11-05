package services

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
)

type ResourceModeler interface {
	GetAll() ([]*model.Resource, error)
	GetOne(int) (*model.Resource, error)
	Create(*types.CreateResourceRequest) error
	Update(int, *types.CreateResourceRequest) (*model.Resource, error)
	Delete(int) error
}

type ResourceService struct {
	roleModel ResourceModeler
}

func (r *ResourceService) GetAll() ([]*model.Resource, error) {
	return r.roleModel.GetAll()
}

func (r *ResourceService) GetOne(id int) (*model.Resource, error) {
	return r.roleModel.GetOne(id)
}

func (r *ResourceService) Create(newResource *types.CreateResourceRequest) error {
	return r.roleModel.Create(newResource)
}

func (r *ResourceService) Update(id int, newResource *types.CreateResourceRequest) (*model.Resource, error) {
	return r.roleModel.Update(id, newResource)
}

func (r *ResourceService) Delete(id int) error {
	return r.roleModel.Delete(id)
}

func NewResourceService(database *db.Storage) *ResourceService {
	roleModel := model.NewResourceModel(database.DB)
	return &ResourceService{
		roleModel: roleModel,
	}
}
