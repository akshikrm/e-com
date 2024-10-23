package user

import "errors"

var UserNotFound = errors.New("not found")
var Unauthorized = errors.New("unauthorized")
