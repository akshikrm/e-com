package main

import (
	db "akshidas/e-com/pkg/db"
	server "akshidas/e-com/pkg/server"
	"akshidas/e-com/pkg/types"
	"database/sql"
	"github.com/joho/godotenv"
	"log"
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
	return nil, nil
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

func NewUserService(db *db.Store) types.UserService {
	return &userService{DB: db.DB}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database := &db.Store{}
	db.Connect(database)

	userService := NewUserService(database)

	server := &server.APIServer{
		Status: "Server is up and running",
		Port:   ":5234",
		User:   userService,
	}
	server.Run()
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
