package service_test

import (
	"context"
	"testing"

	"github.com/go-rel/reltest"
	"github.com/jarcoal/httpmock"
	"github.com/sea-auca/auc-auth/config"
	"github.com/sea-auca/auc-auth/user/repo"
	"github.com/sea-auca/auc-auth/user/service"
	"github.com/stretchr/testify/assert"
)

func createRepos() (*reltest.Repository, service.UserService) {
	trepo := reltest.New()
	ur := repo.NewUserRepository(trepo)
	vlr := repo.NewVerificationRepository(trepo)
	serv := service.NewService(ur, vlr)
	return trepo, serv
}

func TestRegisterUserSuccesful(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.EditbleConfig().Email.Host = "http://localhost"
	config.EditbleConfig().Email.Port = 8002

	httpmock.RegisterResponder("POST", "http://localhost:8002/send/registration",
		httpmock.NewStringResponder(200, ``))

	trp, srv := createRepos()
	trp.ExpectInsert().ForType("service.User")
	trp.ExpectInsert().ForType("service.VerificationLink")
	email := "student_t@auca.kg"
	err := srv.RegisterUser(context.TODO(), email)
	assert.NoError(t, err)
	trp.AssertExpectations(t)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRegisterUserIncorrectEmail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.EditbleConfig().Email.Host = "http://localhost"
	config.EditbleConfig().Email.Port = 8002

	httpmock.RegisterResponder("POST", "http://localhost:8002/send/registration",
		httpmock.NewStringResponder(200, ``))

	trp, srv := createRepos()
	email := "student@auca.kg"
	err := srv.RegisterUser(context.TODO(), email)
	assert.Error(t, err)
	trp.AssertExpectations(t)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestDeactivateAccountSuccessful(t *testing.T) {
	trp, srv := createRepos()
	user := service.NewUser("student_t@auca.kg")
	user.IsActive = true
	trp.ExpectUpdate().For(user).Success()
	err := srv.DeactivateAccount(context.TODO(), user)
	assert.NoError(t, err)
	trp.AssertExpectations(t)
	assert.False(t, user.IsActive)
}
