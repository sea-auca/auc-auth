package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/sea-auca/auc-auth/user/repo"
	"github.com/sea-auca/auc-auth/user/service"

	"github.com/go-rel/rel"
	"github.com/go-rel/reltest"
	"github.com/stretchr/testify/assert"
)

var testUser = service.NewUser("email@domain.com")

func TestCreate(t *testing.T) {
	rep := reltest.New()
	r := repo.NewUserRepository(rep)
	u := *testUser //new user is created
	rep.ExpectInsert().ForType("*service.User")
	_, err := r.Create(context.TODO(), &u)
	assert.Nil(t, err)
	rep.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	rep := reltest.New()
	r := repo.NewUserRepository(rep)
	u := *testUser // suppose the testUser replaces old value
	rep.ExpectUpdate().ForType("*user.User")
	err := r.Update(context.TODO(), &u)
	assert.Nil(t, err)
	rep.AssertExpectations(t)
}

func TestPaginatedView(t *testing.T) {
	rep := reltest.New()
	r := repo.NewUserRepository(rep)
	u := *testUser
	u.CreatedAt = time.Now().Add(time.Minute)
	result := []*service.User{&u, testUser}
	rep.ExpectFindAndCountAll(
		rel.Select().SortDesc("created_at").Limit(5).Offset(1*5),
	).Result(result, len(result))

	assert.NotPanics(t, func() {
		us, cnt, err := r.PaginatedView(context.TODO(), 1, 5)
		assert.Nil(t, err)
		assert.Equal(t, cnt, len(result))
		assert.Equal(t, result, us)
	})
}
