package svcerror

import "errors"

var (
	ErrFailedToDecode      = errors.New("failed to decode request")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrAlreadyExists       = errors.New("already exists")
	ErrNotFound            = errors.New("not found")
	ErrBadRouting          = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrNotAuthorized       = errors.New("user is not authorized to access the resources")
	ErrLoginFailed         = errors.New("failed to login user")
	ErrInvalidToken        = errors.New("invalid JWT token")
	ErrFailedToGenerateJWT = errors.New("failed to generate jwt token")
	ErrFailedToCreateUser  = errors.New("failed to create user")
	ErrFailedToGetUser     = errors.New("failed to get user")
	ErrFailedToListUsers   = errors.New("failed to list users")
	ErrFailedToUpdateUser  = errors.New("failed to update user")
	ErrFailedToDeleteUser  = errors.New("failed to delete user")

	// invoice error
	ErrFailedToCreateInvoice = errors.New("failed to create invoice")
	ErrFailedToDeleteInvoice = errors.New("failed to delete invoice")
	ErrFailedToUpdateInvoice = errors.New("failed to update invoice")
	ErrFailedToGetInvoice    = errors.New("failed to get invoice")
	ErrFailedToListInvoice   = errors.New("failed to list invoice")

	ErrInvalidLoginCreds = errors.New("email or password is wrong")
)
