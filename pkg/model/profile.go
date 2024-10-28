package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Profile struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Pincode     string    `json:"pincode"`
	AddressOne  string    `json:"address_one"`
	AddressTwo  string    `json:"address_two"`
	PhoneNumber string    `json:"phone_number"`
	UserID      int       `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProfileModel struct {
	DB *sql.DB
}

func (p *ProfileModel) GetByUserId(userId int) (*Profile, error) {
	query := `select * from profiles where user_id=$1`
	row := p.DB.QueryRow(query, userId)

	savedProfile := &Profile{}
	if err := row.Scan(
		&savedProfile.ID,
		&savedProfile.UserID,
		&savedProfile.Pincode,
		&savedProfile.AddressOne,
		&savedProfile.AddressTwo,
		&savedProfile.PhoneNumber,
		&savedProfile.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFound
		}
		log.Printf("failed to get profile of user with id: %d due to: %s", userId, err)
		return nil, utils.ServerError
	}
	fmt.Println(*savedProfile)
	return savedProfile, nil
}

func (p *ProfileModel) Create(profile *types.NewProfileRequest) (int, error) {
	query := `insert into
	profiles (first_name,last_name, email, pincode, address_one, address_two, phone_number, user_id, created_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	returning id
	`

	log.Println("Create profile")
	row := p.DB.QueryRow(query,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.AddressOne,
		profile.AddressTwo,
		profile.Pincode,
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

func (p *ProfileModel) UpdateProfileByUserID(userId int, profile *types.UpdateProfileRequest) error {
	query := `update profiles set pincode=$1, address_one=$2, address_two=$3, phone_number=$4 where user_id=$5`

	result, err := p.DB.Exec(query,
		profile.Pincode,
		profile.AddressOne,
		profile.AddressTwo,
		profile.PhoneNumber,
		userId,
	)

	if err != nil {
		log.Printf("failed to update profile %v due to %s", profile, err)
		return utils.ServerError
	}

	if count, _ := result.RowsAffected(); count == 0 {
		log.Printf("updated %d rows", count)
		return utils.NotFound
	}

	return nil
}

func NewProfileModel(DB *sql.DB) *ProfileModel {
	return &ProfileModel{DB: DB}
}
