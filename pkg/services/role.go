package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"database/sql"
)

type RoleModeler interface {
	GetAll() ([]*model.Role, error)
	GetOne(int) (*model.Role, error)
	Create(*types.CreateRoleRequest) error
	Update(int, *types.CreateRoleRequest) (*model.Role, error)
	Delete(int) error
}

type RoleService struct {
	roleModel RoleModeler
}

func (r *RoleService) GetAll() ([]*model.Role, error) {
	return r.roleModel.GetAll()
}

func (r *RoleService) GetOne(id int) (*model.Role, error) {
	return r.roleModel.GetOne(id)
}

func (r *RoleService) Create(newRole *types.CreateRoleRequest) error {
	return r.roleModel.Create(newRole)
}

func (r *RoleService) Update(id int, newRole *types.CreateRoleRequest) (*model.Role, error) {
	return r.roleModel.Update(id, newRole)
}

func (r *RoleService) Delete(id int) error {
	return r.roleModel.Delete(id)
}

func NewRoleService(database *sql.DB) *RoleService {
	roleModel := model.NewRoleModel(database)
	return &RoleService{
		roleModel: roleModel,
	}
}
