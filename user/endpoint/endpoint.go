package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/sea-auca/auc-auth/middleware"
	"github.com/sea-auca/auc-auth/user/service"
	"go.uber.org/zap"
)

type RegistrationRequest struct {
	Email string `json:"email"`
}

type RegistrationResponse struct {
	Errors string `json:"err,omitempty"`
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
