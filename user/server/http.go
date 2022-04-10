package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/transport"
	transp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sea-auca/auc-auth/user/endpoint"
	"github.com/sea-auca/auc-auth/user/service"
	"go.uber.org/zap"
)

type handlers struct {
	regHandler        *transp.Server
	verifyHandler     *transp.Server
	deactivateHandler *transp.Server
	reactivateHandler *transp.Server
	listHandler       *transp.Server
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

	hands.verifyHandler = transp.NewServer(
		endpoint.MakeValidationEndpoint(us, lg),
		decodeVerifyRequest,
		encodeResponse,
		options...,
	)

	hands.deactivateHandler = transp.NewServer(
		endpoint.MakeDeactivationEndpoint(us, lg),
		decodeRequest,
		encodeResponse,
		options...,
	)

	hands.reactivateHandler = transp.NewServer(
		endpoint.MakeReactivationEndpoint(us, lg),
		decodeReactivateRequest,
		encodeResponse,
		options...,
	)

	hands.listHandler = transp.NewServer(
		endpoint.MakeListUsersEndpoint(us, lg),
		decodeRequest,
		encodeResponse,
		options...,
	)

	return hands
}

func RegisterRoutes(us service.UserService, r *mux.Router) {
	hands := MakeHandlers(us)
	r.Path("/v1/register").Methods("POST").Handler(hands.regHandler)
	r.Path("/v1/reactivate").Methods("POST").Handler(hands.reactivateHandler)
	r.Path("/v1/deactivate").Methods("POST").Handler(hands.deactivateHandler)
	r.Path("/v1/verify").Methods("POST").Handler(hands.verifyHandler)
	r.Path("/v1/list").Methods("GET").Handler(hands.listHandler)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeReactivateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.ReactivationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVerifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.ValidationRequest
	var err error
	code := r.URL.Query().Get("code")
	if code == "" {
		return nil, errors.New("malformed request, code has invalid format")
	}
	request.ID, err = uuid.Parse(code)
	if err != nil {
		return nil, errors.New("malformed request, code has invalid format")
	}
	return request, nil
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
