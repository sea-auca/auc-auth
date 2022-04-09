package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/sea-auca/auc-auth/user/service"
)

type RegistrationRequest struct {
	Email string `json:"email"`
}

type RegistrationResponse struct {
	Errors string `json:"err,omitempty"`
}

func MakeRegistrationEndpoint(us service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegistrationRequest)
		response = RegistrationResponse{}
		err = us.RegisterUser(ctx, req.Email)
		if err != nil {
			response = RegistrationResponse{Errors: err.Error()}
		}
		return
	}
}
