package coreErrors

import "errors"

var (
	ErrorIncorrectLoginOrPass = errors.New("incorrect login or password")
	ErrorUserAlreadyLinked    = errors.New("user already linked")
)
