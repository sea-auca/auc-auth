package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	httpT "github.com/go-kit/kit/transport/http"
	"github.com/sea-auca/auc-auth/user/endpoint"
	"github.com/sea-auca/auc-auth/user/service"
)

func NewHTTP(us service.UserService) {
	registrationHandler := httpT.NewServer(
		endpoint.MakeRegistrationEndpoint(us),
		decodeRegisterRequest,
		encodeResponse,
	)
	http.Handle("/user/register", registrationHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
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
