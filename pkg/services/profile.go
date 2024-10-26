package services

import (
	"akshidas/e-com/pkg/types"
)

type ProfileModeler interface {
	Create(types.NewProfileRequest) (int, error)
}

type ProfileService struct {
	model ProfileModeler
}

func (p *ProfileService) Create(profile types.NewProfileRequest) error {
	_, err := p.model.Create(profile)
	if err != nil {
		return err
	}
	return nil
}

func NewProfileService(model ProfileModeler) *ProfileService {
	return &ProfileService{model}
}
