package repository

import (
	"errors"
	"time"
)

type File struct {
	ID         string
	UserID     int64
	Path       string
	SizeBytes  int64
	MimeType   *string
	StorageKey string
	CreatedAt  time.Time
}

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserExists           = errors.New("user already exists")
	ErrExternalUserIDEmpty  = errors.New("externalUserID is empty")
	ErrUserIDMustBePositive = errors.New("userID must be positive")

	ErrFileNotFound     = errors.New("file not found")
	ErrFileExists       = errors.New("file already exists")
	ErrFileTooLarge     = errors.New("file too large")
	ErrFilePathEmpty    = errors.New("path is empty")
	ErrStorageKeyEmpty  = errors.New("storageKey is empty")
	ErrFileSizeNegative = errors.New("sizeBytes is negative")
)
