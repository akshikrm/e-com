package storage

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"fmt"
	"log"
)

type ProfileStorage struct {
	DB *sql.DB
}

func (p *ProfileStorage) GetByUserId(userId int) (*types.Profile, error) {
	query := `select * from profiles where user_id=$1`
	row := p.DB.QueryRow(query, userId)

	savedProfile := &types.Profile{}
	if err := row.Scan(
		&savedProfile.ID,
		&savedProfile.UserID,
		&savedProfile.FirstName,
		&savedProfile.LastName,
		&savedProfile.Email,
		&savedProfile.Pincode,
		&savedProfile.AddressOne,
		&savedProfile.AddressTwo,
		&savedProfile.PhoneNumber,
		&savedProfile.CreatedAt,
		&savedProfile.UpdatedAt,
		&savedProfile.DeletedAt,
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

func (p *ProfileStorage) Create(profile *types.NewProfileRequest) (int, error) {
	query := `insert into
	profiles (first_name,last_name, email, user_id)
	values ($1, $2, $3, $4)
	returning id
	`

	log.Println("Creating profile")
	row := p.DB.QueryRow(query,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.UserID,
	)

	savedProfile := types.Profile{}
	if err := row.Scan(&savedProfile.ID); err != nil {
		log.Printf("failed to scan user after writing %d %s", savedProfile.ID, err)
		return 0, utils.ServerError
	}
	return savedProfile.ID, nil
}

func (p *ProfileStorage) CheckIfUserExists(email string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM profiles WHERE email=$1)"
	row := p.DB.QueryRow(query, email)
	var status bool
	if err := row.Scan(&status); err != nil {
		log.Printf("failed to check if user with %s email exists due to %s", email, err)
		return false
	}
	return status
}

func (p *ProfileStorage) UpdateProfileByUserID(userId int, profile *types.UpdateProfileRequest) error {
	query := `update profiles set pincode=$1, address_one=$2, address_two=$3, phone_number=$4, first_name=$5, last_name=$6, email=$7 where user_id=$8`

	result, err := p.DB.Exec(query,
		profile.Pincode,
		profile.AddressOne,
		profile.AddressTwo,
		profile.PhoneNumber,
		profile.FirstName,
		profile.LastName,
		profile.Email,
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

func NewProfileStorage(DB *sql.DB) *ProfileStorage {
	return &ProfileStorage{DB: DB}
}
