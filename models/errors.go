package models

import "errors"

var (
	ErrNotOwner = errors.New("you do not own this resource")
)
