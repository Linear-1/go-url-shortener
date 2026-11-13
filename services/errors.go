package services


import "errors"

var(

	ErrURLNotFound = errors.New("Short URL not found")
	ErrURLExpired = errors.New("Short URL has expired")

)
