package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"database/sql"
)

type GroupPermissionModeler interface {
	GetAll() ([]*model.GroupPermission, error)
	GetOne(int) (*model.GroupPermission, error)
	Create(*types.CreateNewGroupPermission) error
	Update(int, *types.CreateNewGroupPermission) (*model.GroupPermission, error)
	Delete(int) error
}

type GroupPermissionService struct {
	groupPermissionModel GroupPermissionModeler
}

func (r *GroupPermissionService) GetAll() ([]*model.GroupPermission, error) {
	return r.groupPermissionModel.GetAll()
}

func (r *GroupPermissionService) GetOne(id int) (*model.GroupPermission, error) {
	return r.groupPermissionModel.GetOne(id)
}

func (r *GroupPermissionService) Create(newGroupPermission *types.CreateNewGroupPermission) error {
	return r.groupPermissionModel.Create(newGroupPermission)
}

func (r *GroupPermissionService) Update(id int, newGroupPermission *types.CreateNewGroupPermission) (*model.GroupPermission, error) {
	return r.groupPermissionModel.Update(id, newGroupPermission)
}

func (r *GroupPermissionService) Delete(id int) error {
	return r.groupPermissionModel.Delete(id)
}

func NewGroupPermissionService(database *sql.DB) *GroupPermissionService {
	groupPermissionModel := model.NewGroupPermissionModel(database)
	return &GroupPermissionService{
		groupPermissionModel: groupPermissionModel,
	}
}
