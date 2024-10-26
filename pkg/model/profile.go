package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Profile struct {
	ID          int       `json:"id"`
	Pincode     string    `json:"pincode"`
	AddressOne  string    `json:"address_one"`
	AddressTwo  string    `json:"address_two"`
	PhoneNumber string    `json:"phone_number"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProfileModel struct {
	DB *sql.DB
}

func (p *ProfileModel) Get() {}

func (p *ProfileModel) Create(profile types.NewProfileRequest) (int, error) {
	query := `insert into
	profiles (pincode, address_one, address_two, phone_number, user_id, created_at)
	values ($1, $2, $3, $4, $5, $6)
	returning id
	`

	log.Println("Create profile")
	row := p.DB.QueryRow(query,
		profile.Pincode,
		profile.AddressOne,
		profile.AddressTwo,
		profile.PhoneNumber,
		profile.UserID,
		time.Now().UTC(),
	)

	savedProfile := Profile{}
	if err := row.Scan(&savedProfile.ID); err != nil {
		log.Printf("failed to scan user after writing %d %s", savedProfile.ID, err)
		return 0, utils.ServerError
	}

	return savedProfile.ID, nil
}

func NewProfileModel(DB *sql.DB) *ProfileModel {
	return &ProfileModel{DB: DB}
}
