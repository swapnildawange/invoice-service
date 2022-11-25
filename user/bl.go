package user

import (
	"context"
	"fmt"
	"time"

	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
	"github.com/invoice-service/user/repository"

	"github.com/go-kit/log"
)

//go:generate  mockgen -destination=mocks/bl.mock.go -package=mocks github.com/invoice-service/user BL
type BL interface {
	CreateUser(ctx context.Context, createUserReq spec.CreateUserRequest) (spec.User, error)
	GetUser(ctx context.Context, userId int) (spec.User, error)
	ListUsers(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error)
	DeleteUser(ctx context.Context, deleteUserReq spec.DeleteUserReq) (int, error)
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
		bl.logger.Log("[debug]", fmt.Errorf("failed to create user with email %v error %s", createUserReq.Email, err.Error()))
		return user, err
	}

	user, err = bl.repo.Get(ctx, userId)
	if err != nil {
		bl.logger.Log("[debug]", fmt.Errorf("failed to get user details for %v %v", user.Id, err.Error()))
		return user, err
	}
	user.Email = createUserReq.Email
	bl.logger.Log("[debug]", "Created user successfully", "UserID", user.Id)
	return user, nil
}

func (bl bl) GetUser(ctx context.Context, userId int) (spec.User, error) {
	user, err := bl.repo.Get(ctx, userId)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to get user", "err", err.Error())
		return user, err
	}
	return user, nil
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

func (bl bl) DeleteUser(ctx context.Context, deleteUserReq spec.DeleteUserReq) (int, error) {
	deletedUserId, err := bl.repo.Delete(ctx, deleteUserReq)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to delete user", "err", err.Error())
		return deletedUserId, err
	}
	return deletedUserId, nil
}

func (bl bl) EditUser(ctx context.Context, editUserReq spec.EditUserRequest) (spec.User, error) {
	editedUser, err := bl.repo.Edit(ctx, editUserReq)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to edit user", "err", err.Error())
		return editedUser, err
	}
	return editedUser, nil
}
