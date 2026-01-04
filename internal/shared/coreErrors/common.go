package coreErrors

import (
	"errors"
)

var (
	ErrorNotFound     = errors.New("entity not found")
	ErrorAlreadyExist = errors.New("entity already exist")
)
