package user

import (
	"context"
	"time"

	"invoice_service/model"
	"invoice_service/security"
	"invoice_service/user/repository"

	"github.com/go-kit/kit/log"
)

type BL interface {
	CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	Login(ctx context.Context, loginReq model.LoginRequest) (model.User, string, error)
}

type bl struct {
	logger log.Logger
	repo   repository.Repository
}

func NewBL(logger log.Logger, repo repository.Repository) BL {
	return &bl{
		logger: logger,
		repo:   repo,
	}
}

func (bl *bl) CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error) {
	var user model.User
	// hash password
	hashedPassword, err := security.HashPassword(createUserReq.Password)
	if err != nil {
		bl.logger.Log("CreateUser", "Failed to hash password", err.Error())
	}

	createUserReq.CreatedAt = time.Now()
	createUserReq.UpdatedAt = time.Now()
	createUserReq.Password = hashedPassword

	userId, err := bl.repo.CreateUser(ctx, createUserReq)
	if err != nil {
		bl.logger.Log(err.Error())
		return user, err
	}

	user, err = bl.repo.GetUser(ctx, userId)
	if err != nil {
		bl.logger.Log("Faild to get user details", err)
		return user, err
	}
	user.Email = createUserReq.Email
	bl.logger.Log("User created successfully")
	return user, nil
}

func (bl *bl) Login(ctx context.Context, loginReq model.LoginRequest) (model.User, string, error) {
	var (
		user           model.User
		userId         int
		hashedPassword string
		err            error
		token          string
	)
	// get user details from auth table using email
	userId, hashedPassword, err = bl.repo.GetUserFromAuth(ctx, loginReq.Email)
	if err != nil {
		bl.logger.Log("Failed to login user", err.Error())
		return user, token, err
	}

	err = security.CheckPasswordHash(loginReq.Password, hashedPassword)
	if err != nil {
		bl.logger.Log("Failed to login", err)
		return user, token, err
	}

	// get user details
	user, err = bl.repo.GetUser(ctx, userId)
	if err != nil {
		bl.logger.Log("Faild to get user details", err)
		return user, token, err
	}
	user.Email = loginReq.Email

	// generate jwt token
	token, err = security.GenerateJWT("mysecret", userId, user.Role)
	if err != nil {
		bl.logger.Log("Failed to generate jwt token", err)
		return user, token, err
	}
	return user, token, nil
}

func (bl *bl) ListUsers(ctx context.Context) ([]model.User, error) {
	var (
		users []model.User
		err   error
	)
	users, err = bl.repo.ListUsers(ctx)
	if err != nil {
		bl.logger.Log("Failed to get list of users", err.Error())
		return users, err
	}
	return users, nil
}
