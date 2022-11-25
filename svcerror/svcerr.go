package svcerror

type CustomError interface {
	Error() string
}

type CustomErrString struct {
	s string
}

func NewCustomError(text string) CustomError {
	return &CustomErrString{text}
}

func (e *CustomErrString) Error() string {
	return e.s
}

var (
	ErrFailedToDecode      = NewCustomError("failed to decode request")
	ErrInvalidRequest      = NewCustomError("invalid request")
	ErrAlreadyExists       = NewCustomError("already exists")
	ErrNotFound            = NewCustomError("not found")
	ErrBadRouting          = NewCustomError("inconsistent mapping between route and handler (programmer error)")
	ErrNotAuthorized       = NewCustomError("user is not authorized to access the resources")
	ErrLoginFailed         = NewCustomError("failed to login user")
	ErrInvalidToken        = NewCustomError("invalid JWT token")
	ErrFailedToGenerateJWT = NewCustomError("failed to generate jwt token")
	ErrFailedToCreateUser  = NewCustomError("failed to create user")
	ErrFailedToGetUser     = NewCustomError("failed to get user")
	ErrFailedToListUsers   = NewCustomError("failed to list users")
	ErrFailedToUpdateUser  = NewCustomError("failed to update user")
	ErrFailedToDeleteUser  = NewCustomError("failed to delete user")

	// invoice error
	ErrFailedToCreateInvoice = NewCustomError("failed to create invoice")
	ErrFailedToDeleteInvoice = NewCustomError("failed to delete invoice")
	ErrFailedToUpdateInvoice = NewCustomError("failed to update invoice")
	ErrFailedToGetInvoice    = NewCustomError("failed to get invoice")
	ErrFailedToListInvoice   = NewCustomError("failed to list invoice")
	ErrInvoiceNotFound       = NewCustomError("invoice not found")
	ErrInvalidLoginCreds     = NewCustomError("email or password is wrong")
	ErrSameUserAndAdminId    = NewCustomError("user id and admin id can't be same")

	// jwt
	ErrFailedToGenerateAccessToken  = NewCustomError("failed to generate access token")
	ErrFailedToGenerateRefreshToken = NewCustomError("failed to generate refresh token")

	ErrUserNotFound = NewCustomError("user not found")
)
