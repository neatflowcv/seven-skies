package repository

import "errors"

var (
	ErrWeatherAlreadyExists = errors.New("weather already exists")
	ErrWeatherNotFound      = errors.New("weather not found")
)
