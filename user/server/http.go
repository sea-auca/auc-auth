package server

import (
	"context"
	"encoding/json"
	"net/http"

	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/transport"
	transp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sea-auca/auc-auth/user/endpoint"
	"github.com/sea-auca/auc-auth/user/service"
	"go.uber.org/zap"
)

type handlers struct {
	regHandler *transp.Server
}

type errMsg struct {
	Err string `json:"error"`
}

func MakeHandlers(us service.UserService) handlers {
	lg := zap.L()
	options := []transp.ServerOption{
		transp.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(errMsg{err.Error()})
		}),
		transp.ServerErrorHandler(transport.NewLogErrorHandler(kitzap.NewZapSugarLogger(lg, zap.ErrorLevel))),
	}

	var hands handlers

	hands.regHandler = transp.NewServer(
		endpoint.MakeRegistrationEndpoint(us, lg),
		decodeRegisterRequest,
		encodeResponse,
		options...,
	)

	return hands
}

func RegisterRoutes(us service.UserService, r *mux.Router) {
	hands := MakeHandlers(us)
	r.Path("/send/registration").Methods("POST").Handler(hands.regHandler)
	//r.Path("/send/authentication").Methods("POST").Handler(hands.authHandler)

}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
