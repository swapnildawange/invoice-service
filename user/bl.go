package user

import (
	"context"
	"fmt"
	"time"

	"invoice_service/security"
	"invoice_service/spec"
	"invoice_service/user/repository"

	"github.com/go-kit/kit/log"
)

//go:generate  mockgen -destination=mocks/bl.mock.go -package=mocks invoice_service/user BL
type BL interface {
	CreateUser(ctx context.Context, createUserReq spec.CreateUserRequest) (spec.User, error)
	ListUsers(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error)
	Login(ctx context.Context, loginReq spec.LoginRequest) (spec.User, string, error)
	DeleteUser(ctx context.Context, deleteUserReq spec.DeleteUserReq) (string, error)
	EditUser(ctx context.Context, editUserReq spec.EditUserRequest) (spec.User, error)
}

type bl struct {
	logger log.Logger
	repo   repository.Repository
}

func NewBL(logger log.Logger, repo repository.Repository) BL {
	return bl{
		logger: logger,
		repo:   repo,
	}
}

func (bl bl) CreateUser(ctx context.Context, createUserReq spec.CreateUserRequest) (spec.User, error) {
	var user spec.User
	// hash password
	hashedPassword, err := security.HashPassword(createUserReq.Password)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to hash password", "err", err.Error())
		return user, err
	}

	createUserReq.CreatedAt = time.Now()
	createUserReq.UpdatedAt = time.Now()
	createUserReq.Password = hashedPassword

	userId, err := bl.repo.Create(ctx, createUserReq)
	if err != nil {
		bl.logger.Log("[debug]", fmt.Errorf("failed to create user with email %v error %w", createUserReq.Email, err))
		return user, err
	}

	user, err = bl.repo.Get(ctx, userId)
	if err != nil {
		bl.logger.Log("[debug]", fmt.Errorf("failed to get user details for %v %w", user.Id, err))
		return user, err
	}
	user.Email = createUserReq.Email
	bl.logger.Log("[debug]", "Created user successfully", "UserID", user.Id)
	return user, nil
}

func (bl bl) Login(ctx context.Context, loginReq spec.LoginRequest) (spec.User, string, error) {
	var (
		user           spec.User
		userId         int
		hashedPassword string
		err            error
		token          string
	)

	// get user details from auth table using email
	userId, hashedPassword, err = bl.repo.GetUserFromAuth(ctx, loginReq.Email)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to login user", "err", err.Error())
		return user, token, err
	}

	err = security.CheckPasswordHash(loginReq.Password, hashedPassword)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to login", "err", err.Error())
		return user, token, err
	}

	// get user details
	user, err = bl.repo.Get(ctx, userId)
	if err != nil {
		bl.logger.Log("[debug]", "Faild to get user details", "err", err.Error())
		return user, token, err
	}
	user.Email = loginReq.Email

	// generate jwt token
	token, err = security.GenerateJWT("mysecret", userId, user.Role)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to generate jwt token", "err", err.Error())
		return user, token, err
	}
	return user, token, nil
}

func (bl bl) ListUsers(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error) {
	var (
		users []spec.User
		err   error
	)
	users, err = bl.repo.List(ctx, listUserFilter)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to get list of users", "err", err.Error())
		return users, err
	}
	return users, nil
}

func (bl bl) DeleteUser(ctx context.Context, deleteUserReq spec.DeleteUserReq) (string, error) {
	err := bl.repo.Delete(ctx, deleteUserReq)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to delete user", "err", err.Error())
		return "", fmt.Errorf("failed to delete user %w", err)
	}
	return "Deleted user Successfully", nil
}

func (bl bl) EditUser(ctx context.Context, editUserReq spec.EditUserRequest) (spec.User, error) {
	editedUser, err := bl.repo.Edit(ctx, editUserReq)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to edit user", "err", err.Error())
		return editedUser, fmt.Errorf("failed to edit user %w", err)
	}
	return editedUser, nil
}
