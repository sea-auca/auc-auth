package user

import (
	"context"

	"github.com/go-rel/changeset"
	"github.com/go-rel/changeset/params"
	r "github.com/go-rel/rel"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type postgresRepository struct {
	logger *zap.SugaredLogger
	repo   r.Repository
}

func (repo postgresRepository) ChangeUser(user interface{}, p params.Params) *changeset.Changeset {
	logger := zap.S() // get global logger
	ch := changeset.Cast(user, p, []string{"fullname", "password", "verified"})
	if ch.Get("password") != nil { // the password branch checks if password was changed and updates the hash
		changeset.ValidateMin(ch, "password", 9)
		val, ok := ch.Get("password").([]byte)
		if !ok {
			logger.Errorw("Could not convert password to byte slice", "password value", ch.Get("password")) // do not log password?
			changeset.AddError(ch, "password", "could not convert to byte slice")
			goto AfterPassword
		}
		hash, err := bcrypt.GenerateFromPassword(val, 15)
		if err != nil {
			logger.Errorw("Failed to hash password", "error", err)
			changeset.AddError(ch, "password", "failed to hash password")
			goto AfterPassword
		}
		changeset.PutChange(ch, "hash", hash) // update the hash with new password
	}
AfterPassword:
	changeset.ValidateRequired(ch, []string{"fullname", "verified"})
	changeset.ValidateMin(ch, "fullname", 4)
	changeset.ValidatePattern(ch, "email", aucaEmail)
	return ch
}

func NewUserRepository(repo r.Repository) UserRepository {
	return postgresRepository{logger: zap.S(), repo: repo}
}

func (rp postgresRepository) Create(ctx context.Context, u *User) (*User, error) {
	err := rp.repo.Insert(ctx, u)
	return u, err
}

func (rp postgresRepository) Update(ctx context.Context, u *User) error {
	return rp.repo.Update(ctx, u)
}

func (rp postgresRepository) RelUpdate(ctx context.Context, u *User, ch ...r.Mutator) error {
	return rp.repo.Update(ctx, u, ch...)
}

func (rp postgresRepository) GetByID(ctx context.Context, id UUID) (*User, error) {
	var u *User
	err := rp.repo.Find(ctx, u, r.Where(r.Eq("uuid", id)))
	return u, err
}

func (rp postgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u *User
	err := rp.repo.Find(ctx, u, r.Where(r.Eq("email", email)))
	return u, err
}

func (rp postgresRepository) PaginatedView(ctx context.Context, page, pageSize int) ([]*User, int, error) {
	var users []*User
	cnt, err := rp.repo.FindAndCountAll(ctx, &users, r.Select().SortDesc("created_at").Limit(pageSize).Offset(page*pageSize))
	return users, cnt, err
}
