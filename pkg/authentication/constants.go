package authentication

import "errors"

var (
	AuthErrInvalidToken     = errors.New("invalid token")
	AuthErrExpiredToken     = errors.New("expired token")
	AuthErrMalformedToken   = errors.New("malformed token")
	AuthErrSigningMethod    = errors.New("unexpected signing method")
	AuthErrInvalidSignature = errors.New("invalid signature")
	AuthErrInvalidIssuer    = errors.New("invalid issuer")
)
