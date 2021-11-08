package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type registrationRequest struct {
	Email string `json:"email"`
}

type registrationResponse struct {
	Errors string `json:"err,omitempty"`
}

func makeRegistrationEndpoint(us UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registrationRequest)
		response = registrationResponse{}
		err = us.RegisterUser(ctx, req.Email)
		if err != nil {
			response = registrationResponse{Errors: err.Error()}
		}
		return
	}
}
