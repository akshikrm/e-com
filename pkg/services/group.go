package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"database/sql"
)

type GroupModeler interface {
	GetAll() ([]*model.Group, error)
	GetOne(int) (*model.Group, error)
	Create(*types.CreateNewGroup) error
	Update(int, *types.CreateNewGroup) (*model.Group, error)
	Delete(int) error
}

type GroupService struct {
	groupModel GroupModeler
}

func (r *GroupService) GetAll() ([]*model.Group, error) {
	return r.groupModel.GetAll()
}

func (r *GroupService) GetOne(id int) (*model.Group, error) {
	return r.groupModel.GetOne(id)
}

func (r *GroupService) Create(newGroup *types.CreateNewGroup) error {
	return r.groupModel.Create(newGroup)
}

func (r *GroupService) Update(id int, newGroup *types.CreateNewGroup) (*model.Group, error) {
	return r.groupModel.Update(id, newGroup)
}

func (r *GroupService) Delete(id int) error {
	return r.groupModel.Delete(id)
}

func NewGroupService(database *sql.DB) *GroupService {
	groupModel := model.NewGroupModel(database)
	return &GroupService{
		groupModel: groupModel,
	}
}
