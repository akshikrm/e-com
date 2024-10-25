package model

import (
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
	PhoneNubmer string    `json:"phone_number"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type NewProfileRequest struct {
	Pincode     string `json:"pincode"`
	AddressOne  string `json:"address_one"`
	AddressTwo  string `json:"address_two"`
	PhoneNubmer string `json:"phone_number"`
	UserID      int    `json:"user_id"`
}

type ProfileModel struct {
	DB *sql.DB
}

func (p *ProfileModel) Get() {}

func (p *ProfileModel) Create(profile NewProfileRequest) (int, error) {
	query := `insert into
	profiles (pincode, address_one, address_two, phone_number, user_id, created_at)
	values($1, $2, $3, $3, $4, $5)
	`

	row := p.DB.QueryRow(query,
		profile.Pincode,
		profile.AddressOne,
		profile.AddressTwo,
		profile.PhoneNubmer,
		profile.UserID,
		time.Now().UTC(),
	)

	log.Printf("Create profile %v", profile)

	savedProfile := Profile{}
	if err := row.Scan(&savedProfile.ID); err != nil {
		log.Printf("failed to scan user after writing %d %s", savedProfile.ID, err)
		return 0, utils.ServerError
	}

	return savedProfile.ID, nil
}
