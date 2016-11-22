package main

import "errors"

var (
	ErrNoValidPazuzufile = errors.New("No valid Pazuzufile provided")
	ErrNotImplemented    = errors.New("Feature not implemented yet :(")
	ErrTooManyParameters = errors.New("Too many parameters provided")
	ErrStopIteration     = errors.New("It's not an real error, sorry!")
	ErrNotFound          = errors.New("Not found")
)
