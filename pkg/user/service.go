package user

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/types"
	"database/sql"
	"errors"
	"time"
)

var UserNotFound = errors.New("not found")

func NewUserService(db *db.PostgresStore) *UserService {
	return &UserService{DB: db.DB}
}

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

func (u *UserService) Create(user *types.User) error {
	sqlQuery := `insert into 
	users (first_name, last_name, password, email, created_at)
	values($1, $2, $3, $4, $5)`

	_, err := u.DB.Exec(sqlQuery,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Email,
		time.Now().UTC(),
	)

	return err
}

func (u *UserService) Update(user types.User) error {
	return nil
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
