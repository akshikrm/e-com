package utils

import "errors"

var NotFound = errors.New("not found")
var Failed = errors.New("failed")
var Unauthorized = errors.New("unauthorized")
