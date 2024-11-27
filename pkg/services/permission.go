package services

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/storage"
	"akshidas/e-com/pkg/types"
)

type PermissionModeler interface {
	GetAll() ([]*types.Permission, error)
	GetOne(int) (*types.Permission, error)
	Create(*types.CreateNewPermission) error
	Update(int, *types.CreateNewPermission) (*types.Permission, error)
	Delete(int) error
}

type PermissionService struct {
	permissionModel PermissionModeler
}

func (r *PermissionService) GetAll() ([]*types.Permission, error) {
	return r.permissionModel.GetAll()
}

func (r *PermissionService) GetOne(id int) (*types.Permission, error) {
	return r.permissionModel.GetOne(id)
}

func (r *PermissionService) Create(newPermission *types.CreateNewPermission) error {
	return r.permissionModel.Create(newPermission)
}

func (r *PermissionService) Update(id int, newPermission *types.CreateNewPermission) (*types.Permission, error) {
	return r.permissionModel.Update(id, newPermission)
}

func (r *PermissionService) Delete(id int) error {
	return r.permissionModel.Delete(id)
}

func NewPermissionService(database *db.Storage) *PermissionService {
	permissionModel := storage.NewPermissionStorage(database.DB)
	return &PermissionService{
		permissionModel: permissionModel,
	}
}
