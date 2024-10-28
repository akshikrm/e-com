package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
)

type ProfileModeler interface {
	GetByUserId(int) (*model.Profile, error)
	Create(*types.NewProfileRequest) (int, error)
	UpdateProfileByUserID(int, *types.UpdateProfileRequest) error
}

type ProfileService struct {
	db    *sql.DB
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

func (p *ProfileService) Update(userId int, profile *types.UpdateProfileRequest) (*model.Profile, error) {
	err := p.model.UpdateProfileByUserID(userId, profile)
	if err != nil {
		return nil, err
	}
	return p.GetByUserId(userId)
}

func (p *ProfileService) createAndReturnProfileFromUserId(userID int) (*model.Profile, error) {
	err := p.Create(&types.NewProfileRequest{UserID: userID})
	if err != nil {
		return nil, err
	}

	return p.GetByUserId(userID)
}

func NewProfileService(database *sql.DB) *ProfileService {
	profileModel := model.NewProfileModel(database)
	return &ProfileService{model: profileModel, db: database}
}
