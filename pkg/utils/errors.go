package utils

import "errors"

var NotFound = errors.New("not found")
var ServerError = errors.New("server error")
var Conflict = errors.New("conflict")
var Unauthorized = errors.New("unauthorized")
var InvalidRequest = errors.New("invalid request")
