package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sea-auca/auc-auth/config"
	"go.uber.org/zap"
)

type registrationReq struct {
	To   string `json:"to"`
	Link string `json:"link"`
}

type registrationResp struct {
	Error interface{} `json:"error"`
}

func (s service) sendVereficationEmail(ctx context.Context, code uuid.UUID, email string) error {
	conf := config.Config().Email
	servConf := config.Config().Service
	var reqBody registrationReq
	reqBody.Link = servConf.VerificationPrefix + "/verify?code=" + code.String() + "&action=i"
	reqBody.To = email
	endpoint := fmt.Sprintf("%s:%d/send/registration", conf.Host, conf.Port)
	marhshaledReq, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(marhshaledReq))
	if err != nil {
		return err
	}
	respCont, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if respCont.StatusCode == 200 {
		return nil
	}
	var resp registrationResp
	err = json.NewDecoder(respCont.Body).Decode(&resp)
	if err != nil {
		s.lg.Error("Failed to unmarshal mailer response", zap.Error(err))
		return err
	}
	if resp.Error != nil {
		return errors.New(resp.Error.(string))
	}
	return nil
}
