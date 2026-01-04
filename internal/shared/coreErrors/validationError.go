package coreErrors

import "fmt"

type ValidationError struct {
	Msg    string   `json:"error"`
	Entity string   `json:"entity"`
	Fields []string `json:"fields"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: entity: %s, msg: %s, fields: %v", e.Entity, e.Msg, e.Fields)
}
