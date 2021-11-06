package user_test

import (
	"context"
	"sea/auth/user"
	"testing"
	"time"

	"github.com/go-rel/rel"
	"github.com/go-rel/reltest"
	"github.com/stretchr/testify/assert"
)

var testUser = user.NewUser("email@domain.com")

func TestCreate(t *testing.T) {
	repo := reltest.New()
	r := user.NewUserRepository(repo)
	u := *testUser //new user is created
	repo.ExpectInsert().ForType("*user.User")
	_, err := r.Create(context.TODO(), &u)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	repo := reltest.New()
	r := user.NewUserRepository(repo)
	u := *testUser // suppose the testUser replaces old value
	repo.ExpectUpdate().ForType("*user.User")
	err := r.Update(context.TODO(), &u)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestPaginatedView(t *testing.T) {
	repo := reltest.New()
	r := user.NewUserRepository(repo)
	u := *testUser
	u.CreatedAt = time.Now().Add(time.Minute)
	result := []*user.User{&u, testUser}
	repo.ExpectFindAndCountAll(
		rel.Select().SortDesc("created_at").Limit(5).Offset(1*5),
	).Result(result, len(result))

	assert.NotPanics(t, func() {
		us, cnt, err := r.PaginatedView(context.TODO(), 1, 5)
		assert.Nil(t, err)
		assert.Equal(t, cnt, len(result))
		assert.Equal(t, result, us)
	})
}
