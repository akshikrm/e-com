package types

import "time"

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserService interface {
	Get() ([]*User, error)
	GetOne(id int) (*User, error)
	Create(user *User) error
	Update(user *User) (*User, error)
	Delete(id int) error
}

// type UserDatabase interface {
// 	Get() (*[]User, error)
// 	GetOne(id int) (*User, error)
// 	Create(user User) error
// 	Update(user User) error
// 	Delete(id int) error
// }
