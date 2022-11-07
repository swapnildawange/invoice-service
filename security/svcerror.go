package security

import "errors"

var (
	InvalidLoginErr  = errors.New("Username or Password does not match. Authentication failed.")
	NotAuthorizedErr = errors.New("User is not authorized to access this URL")
)
