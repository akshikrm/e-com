package types

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewProfileRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Pincode     string `json:"pincode"`
	AddressOne  string `json:"address_one"`
	AddressTwo  string `json:"address_two"`
	PhoneNumber string `json:"phone_number"`
	UserID      int    `json:"user_id"`
}

type UpdateProfileRequest struct {
	Pincode     string `json:"pincode"`
	AddressOne  string `json:"address_one"`
	AddressTwo  string `json:"address_two"`
	PhoneNumber string `json:"phone_number"`
}
