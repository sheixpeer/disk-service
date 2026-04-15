package repository

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")

	ErrFileNotFound = errors.New("file not found")
	ErrFileExists   = errors.New("file already exists")
	ErrFileTooLarge = errors.New("file too large")
)
