package usecase

import (
	"context"
	"database/sql"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/common"
	"github.com/lengocson131002/go-clean/pkg/logger"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
	"github.com/lengocson131002/go-clean/pkg/validation"
	"github.com/lengocson131002/go-clean/usecase/data"
	"golang.org/x/crypto/bcrypt"
)

type createUserHandler struct {
	Log            logger.Logger
	Validator      validation.Validator
	UserRepository data.UserRepository
}

func NewCreateUserHandler(log logger.Logger, validator validation.Validator, userRepository data.UserRepository) domain.CreateUserHandler {
	return &createUserHandler{
		Log:            log,
		Validator:      validator,
		UserRepository: userRepository,
	}
}

func (c *createUserHandler) Handle(ctx context.Context, request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	// begin traction

	err := c.Validator.Validate(request)
	if err != nil {
		c.Log.Warn("Invalid request body : %+v", err)
		return nil, common.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(ctx, request.ID)
	if err != nil {
		c.Log.Warn("Failed count user from database : %+v", err)
		return nil, common.ErrInternalServer
	}

	if total > 0 {
		c.Log.Warn("User already exists")
		return nil, domain.ErrorAccountExisted
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warn("Failed to generate bcrype hash : %+v", err)
		return nil, common.ErrInternalServer
	}

	user := &domain.User{
		ID:       request.ID,
		Password: string(password),
		Name:     request.Name,
	}

	err = c.UserRepository.WithinTransactionOptions(ctx, func(ctx context.Context) error {
		return c.UserRepository.CreateUser(ctx, user)
	}, &sql.TxOptions{
		Isolation: sql.IsolationLevel(2),
		ReadOnly:  false,
	})

	if err != nil {
		c.Log.Warn("Failed create user to database : %+v", err)
		return nil, common.ErrInternalServer
	}

	res := &domain.CreateUserResponse{}
	err = mapper.BindingStruct(user, res)
	return res, err
}
