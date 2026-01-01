package coreErrors

import "errors"

var (
	ErrorIncorrectLoginOrPass = errors.New("incorrect login or password")
	ErrorInvalidPassword      = errors.New("user password is invalid")
	ErrorLoginIsTooSmall      = errors.New("login is too small")
	ErrorIncorrectRole        = errors.New("role is incorrect")

	ErrorLoginOccupied = errors.New("user with this login already exist")

	ErrorUserAlreadyLinked = errors.New("user already linked")
)
