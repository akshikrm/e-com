package user

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/types"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type UserService struct {
	DB *sql.DB
}

func (u *UserService) Get() ([]*types.User, error) {
	query := `select * from users;`

	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}

	users := []*types.User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserService) GetOne(id int) (*types.User, error) {
	query := `select * from users where id=$1`
	rows, err := u.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, UserNotFound
}

func (u *UserService) Create(user *types.User) (string, error) {
	query := `insert into 
	users (first_name, last_name, password, email, created_at)
	values($1, $2, $3, $4, $5)
	returning id
	`

	hashedPassword, err := hashPassword([]byte(user.Password))
	if err != nil {
		return "", err
	}

	row := u.DB.QueryRow(query,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Email,
		time.Now().UTC(),
	)

	log.Printf("Created user %v", user)

	savedUser := &types.User{}
	if err := row.Scan(&savedUser.ID); err != nil {
		log.Printf("failed to scan user after saving %v", err)
		return "", err
	}
	return createJwt(savedUser)
}

func (u *UserService) Update(user *types.User) (*types.User, error) {
	query := `update users set first_name=$1, last_name=$2, email=$3 where id=$4`
	result, err := u.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.ID)

	if err != nil {
		log.Printf("failed to update user %v due to %s", user, err)
		return nil, fmt.Errorf("failed to update")
	}

	if count, _ := result.RowsAffected(); count == 0 {
		log.Printf("updated %d rows", count)
		return nil, UserNotFound
	}

	return u.GetOne(user.ID)
}

func (u *UserService) Delete(id int) error {
	query := "delete from users where id=$1"
	result, err := u.DB.Exec(query, id)
	if count, _ := result.RowsAffected(); count == 0 {
		return UserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func NewUserService(db *db.PostgresStore) *UserService {
	return &UserService{DB: db.DB}
}
