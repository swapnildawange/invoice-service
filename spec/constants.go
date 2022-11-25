package spec

const (
	PageSize = 5
	Timeout  = 5
)

// payment status
type PaymentStatus int

const (
	Initiate   PaymentStatus = 1
	InProgress PaymentStatus = 2
	Success    PaymentStatus = 3
	Failed     PaymentStatus = 4
	Retry      PaymentStatus = 5
)

type Role int

const (
	RoleAdmin Role = 1
	RoleUser  Role = 2
)

const (
	CreateUserRequestPath  = "/user"
	ListUsersRequestPath   = "/users"
	EditUserRequestPath    = "/user/{id}"
	GetUserRequestPath     = "/user/{id}"
	DeleteUserRequestPath  = "/user"
	LoginRequestPath       = "/login"
	GenerateJWTRequestPath = "/generate_token"
)
