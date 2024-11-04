package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Password  string    `json:"-"`
	Role      string    `json:"role_code"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Get() ([]*User, error) {
	query := `select * from users;`

	rows, err := m.DB.Query(query)
	if err != nil {
		log.Printf("failed to retrieve user %s", err)
		return nil, utils.ServerError
	}

	users := []*User{}
	for rows.Next() {
		user, err := ScanRows(rows)
		if err != nil {
			return nil, utils.ServerError
		}
		users = append(users, user)
	}

	return users, nil
}

func (m *UserModel) GetPasswordByEmail(email string) (*User, error) {
	query := `select user_id,password from users inner join profiles on users.id = profiles.user_id where email=$1;`
	row := m.DB.QueryRow(query, email)

	user := User{}
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("profile with email: %s not found", email)
			return nil, utils.NotFound
		}
		log.Printf("failed to retrieve for email: %s due to error:%s", email, err)
		return nil, utils.ServerError
	}
	fmt.Println(user)
	return &user, nil
}

func (m *UserModel) GetOne(id int) (*User, error) {
	query := `select * from users where id=$1`
	row := m.DB.QueryRow(query, id)

	user, err := ScanRow(row)
	if err != nil {
		log.Printf("user with id %d not found due to %s", id, err)
		return nil, utils.NotFound
	}

	return user, nil
}

func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	query := `select * from users where email=$1`
	row := m.DB.QueryRow(query, email)

	user, err := ScanRow(row)
	if err != nil {
		log.Printf("user with email %s not found due to %s", email, err)
		return nil, utils.NotFound
	}

	return user, nil

}

func (m *UserModel) Create(user types.CreateUserRequest) (int, error) {
	query := `insert into 
	users (password, role_code)
	values($1, $2)
	returning id
	`
	role := "user"
	if user.Role != "" {
		role = user.Role
	}
	row := m.DB.QueryRow(query,
		user.Password,
		role,
	)
	log.Printf("Created user %v", user)

	savedUser := User{}
	if err := row.Scan(&savedUser.ID); err != nil {
		log.Printf("failed to scan user after writing %d %s", savedUser.ID, err)
		return 0, utils.ServerError
	}

	return savedUser.ID, nil
}

func (m *UserModel) Update(id int, user types.UpdateUserRequest) error {
	query := `update users set first_name=$1, last_name=$2, email=$3 where id=$4`
	result, err := m.DB.Exec(query, user.FirstName, user.LastName, user.Email, id)

	if err != nil {
		log.Printf("failed to update user %v due to %s", user, err)
		return utils.ServerError
	}

	if count, _ := result.RowsAffected(); count == 0 {
		log.Printf("updated %d rows", count)
		return utils.NotFound
	}

	return nil
}

func (m *UserModel) Delete(id int) error {
	query := "delete from users where id=$1"
	result, err := m.DB.Exec(query, id)

	if err != nil {
		log.Printf("failed to delete %d due to %s", id, err)
		return utils.ServerError
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return utils.NotFound
	}

	return nil
}

func ScanRows(rows *sql.Rows) (*User, error) {
	user := &User{}
	err := rows.Scan(
		&user.ID,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		log.Printf("scan into user failed due to %s", err)
	}

	return user, err
}

func ScanRow(row *sql.Row) (*User, error) {
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Password,
		&user.CreatedAt,
	)

	return user, err
}

func NewUserModel(store *sql.DB) *UserModel {
	return &UserModel{
		DB: store,
	}
}
