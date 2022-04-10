package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/google/uuid"
	"github.com/sea-auca/auc-auth/middleware"
	"github.com/sea-auca/auc-auth/user/service"
	"go.uber.org/zap"
)

type RegistrationRequest struct {
	Email string `json:"email"`
}

type RegistrationResponse struct {
	Errors string `json:"error,omitempty"`
}

func MakeRegistrationEndpoint(us service.UserService, lg *zap.Logger) endpoint.Endpoint {
	logger := kitzap.NewZapSugarLogger(lg, zap.InfoLevel)
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegistrationRequest)
		response = RegistrationResponse{}
		err = us.RegisterUser(ctx, req.Email)
		if err != nil {
			response = RegistrationResponse{Errors: err.Error()}
		}
		return
	}
	e = middleware.LoggingMiddleware(logger)(e)
	return e
}

type DeactivationRequest struct {
	UserID uuid.UUID
}

func MakeDeactivationEndpoint(us service.UserService, lg *zap.Logger) endpoint.Endpoint {
	logger := kitzap.NewZapSugarLogger(lg, zap.InfoLevel)
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeactivationRequest)
		user, err := us.GetUserByID(ctx, req.UserID)
		if err != nil {
			return
		}
		err = us.DeactivateAccount(ctx, user)
		return
	}
	e = middleware.LoggingMiddleware(logger)(e)
	return e
}

type ReactivationRequest struct {
	Email string `json:"email"`
}

func MakeReactivationEndpoint(us service.UserService, lg *zap.Logger) endpoint.Endpoint {
	logger := kitzap.NewZapSugarLogger(lg, zap.InfoLevel)
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ReactivationRequest)
		err = us.ReactivateAccount(ctx, req.Email)
		return
	}
	e = middleware.LoggingMiddleware(logger)(e)
	return e
}

type ValidationRequest struct {
	ID uuid.UUID
}

func MakeValidationEndpoint(us service.UserService, lg *zap.Logger) endpoint.Endpoint {
	logger := kitzap.NewZapSugarLogger(lg, zap.InfoLevel)
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ValidationRequest)
		err = us.ValidateUser(ctx, req.ID)
		return
	}
	e = middleware.LoggingMiddleware(logger)(e)
	return e
}

type ListUsersResponse struct {
	Users []*service.User `json:"users"`
}

func MakeListUsersEndpoint(us service.UserService, lg *zap.Logger) endpoint.Endpoint {
	logger := kitzap.NewZapSugarLogger(lg, zap.InfoLevel)
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, _, err := us.ListUsers(ctx, 1, 100000)
		response = ListUsersResponse{users}
		return
	}
	e = middleware.LoggingMiddleware(logger)(e)
	return e
}

func MakeEcho(us service.UserService, lg *zap.Logger) endpoint.Endpoint {
	logger := kitzap.NewZapSugarLogger(lg, zap.InfoLevel)
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("hi")
		return
	}
	e = middleware.LoggingMiddleware(logger)(e)
	return e
}
