package ginjwt

import "errors"

var (
	// ErrInvalidSigningKey is the error returned when a token can not be verified because the signing key in invalid
	ErrInvalidSigningKey = errors.New("unable to find appropriate signing key")

	// ErrInvalidAudience is the error returned when the audience of the token isn't what we expect
	ErrInvalidAudience = errors.New("invalid JWT audience")

	// ErrInvalidIssuer is the error returned when the issuer of the token isn't what we expect
	ErrInvalidIssuer = errors.New("invalid JWT issuer")
)
