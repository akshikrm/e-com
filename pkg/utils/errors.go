package utils

import "errors"

var NotFound = errors.New("not found")
var ServerError = errors.New("server error")
var Unauthorized = errors.New("unauthorized")
