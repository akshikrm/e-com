package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"database/sql"
)

type PermissionModeler interface {
	GetAll() ([]*model.Permission, error)
	GetOne(int) (*model.Permission, error)
	Create(*types.CreateNewPermission) error
	Update(int, *types.CreateNewPermission) (*model.Permission, error)
	Delete(int) error
}

type PermissionService struct {
	permissionModel PermissionModeler
}

func (r *PermissionService) GetAll() ([]*model.Permission, error) {
	return r.permissionModel.GetAll()
}

func (r *PermissionService) GetOne(id int) (*model.Permission, error) {
	return r.permissionModel.GetOne(id)
}

func (r *PermissionService) Create(newPermission *types.CreateNewPermission) error {
	return r.permissionModel.Create(newPermission)
}

func (r *PermissionService) Update(id int, newPermission *types.CreateNewPermission) (*model.Permission, error) {
	return r.permissionModel.Update(id, newPermission)
}

func (r *PermissionService) Delete(id int) error {
	return r.permissionModel.Delete(id)
}

func NewPermissionService(database *sql.DB) *PermissionService {
	permissionModel := model.NewPermissionModel(database)
	return &PermissionService{
		permissionModel: permissionModel,
	}
}
