package service

import "errors"

var ErrEmailExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
