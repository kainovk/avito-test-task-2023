package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user exists")
	ErrUserNotExists = errors.New("user not exists")

	ErrSegmentNotFound  = errors.New("segment not found")
	ErrSegmentExists    = errors.New("segment exists")
	ErrSegmentNotExists = errors.New("segment not exists")

	ErrUserAlreadyHaveSegment = errors.New("user already have segment")
)
