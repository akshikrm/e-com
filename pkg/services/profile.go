package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
)

type ProfileModeler interface {
	GetByUserId(int) (*model.Profile, error)
	Create(*types.NewProfileRequest) (int, error)
}

type ProfileService struct {
	model ProfileModeler
}

func (p *ProfileService) GetByUserId(userID int) (*model.Profile, error) {
	profile, err := p.model.GetByUserId(userID)

	if err == nil {
		return profile, nil
	}

	if err != utils.NotFound {
		return nil, err
	}
	return p.createAndReturnProfileFromUserId(userID)
}

func (p *ProfileService) Create(profile *types.NewProfileRequest) error {
	_, err := p.model.Create(profile)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProfileService) createAndReturnProfileFromUserId(userID int) (*model.Profile, error) {
	err := p.Create(&types.NewProfileRequest{UserID: userID})
	if err != nil {
		return nil, err
	}

	return p.GetByUserId(userID)
}

func NewProfileService(model ProfileModeler) *ProfileService {
	return &ProfileService{model}
}
