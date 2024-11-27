package services

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
)

type RoleModeler interface {
	GetAll() ([]*types.Role, error)
	GetOne(int) (*types.Role, error)
	Create(*types.CreateRoleRequest) error
	Update(int, *types.CreateRoleRequest) (*types.Role, error)
	Delete(int) error
}

type RoleService struct {
	roleModel RoleModeler
}

func (r *RoleService) GetAll() ([]*types.Role, error) {
	return r.roleModel.GetAll()
}

func (r *RoleService) GetOne(id int) (*types.Role, error) {
	return r.roleModel.GetOne(id)
}

func (r *RoleService) Create(newRole *types.CreateRoleRequest) error {
	return r.roleModel.Create(newRole)
}

func (r *RoleService) Update(id int, newRole *types.CreateRoleRequest) (*types.Role, error) {
	return r.roleModel.Update(id, newRole)
}

func (r *RoleService) Delete(id int) error {
	return r.roleModel.Delete(id)
}

func NewRoleService(database *db.Storage) *RoleService {
	roleModel := model.NewRoleStorage(database.DB)
	return &RoleService{
		roleModel: roleModel,
	}
}
