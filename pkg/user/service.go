package user

import (
	"akshidas/e-com/pkg/types"
	"database/sql"
	"fmt"
	"time"
)

type userService struct {
	DB *sql.DB
}

func (u *userService) Get() ([]*types.User, error) {
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

func (u *userService) GetOne(id int) (*types.User, error) {
	query := `select * from users where id=$1`
	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user with id %d not found", id)
}

func (u *userService) Create(user *types.User) error {
	sqlQuery := `insert into 
	users (first_name, last_name, password, email, created_at)
	values($1, $2, $3, $4, $5)`

	_, err := u.DB.Query(sqlQuery,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Email,
		time.Now().UTC(),
	)

	return err
}

func (u *userService) Update(user types.User) error {
	return nil
}

func (u *userService) Delete(id int) error {
	return nil
}

func scanIntoUser(rows *sql.Rows) (*types.User, error) {
	user := &types.User{}
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	return user, err
}
