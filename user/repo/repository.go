package repo

import (
	"context"
	"errors"
	"time"

	"github.com/go-rel/changeset"
	"github.com/go-rel/changeset/params"
	r "github.com/go-rel/rel"
	"github.com/google/uuid"
	"github.com/sea-auca/auc-auth/user/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type postgresRepository struct {
	logger *zap.SugaredLogger
	repo   r.Repository
}

func ChangeUser(user interface{}, p params.Params) *changeset.Changeset {
	logger := zap.S() // get global logger
	ch := changeset.Cast(user, p, []string{"fullname", "password", "email", "verified"})
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
	if ch.Values()["uuid"].(string) == "" {
		newUser := service.NewUser(ch.Get("email").(string))
		changeset.PutChange(ch, "uuid", newUser.ID)
	}
	changeset.ValidateRequired(ch, []string{"email", "active", "uuid"})
	return ch
}

func NewUserRepository(repo r.Repository) service.UserRepository {
	return postgresRepository{logger: zap.S(), repo: repo}
}

func (rp postgresRepository) Create(ctx context.Context, u *service.User) (*service.User, error) {
	err := rp.repo.Insert(ctx, u)
	return u, err
}

func (rp postgresRepository) Update(ctx context.Context, u *service.User) error {
	return rp.repo.Update(ctx, u)
}

func (rp postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*service.User, error) {
	var u service.User
	err := rp.repo.Find(ctx, &u, r.Where(r.Eq("id", id)))
	return &u, err
}

func (rp postgresRepository) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	var u *service.User
	err := rp.repo.Find(ctx, u, r.Where(r.Eq("email", email)))
	return u, err
}

func (rp postgresRepository) PaginatedView(ctx context.Context, page, pageSize int) ([]*service.User, int, error) {
	var users []*service.User
	cnt, err := rp.repo.FindAndCountAll(ctx, &users, r.Select().SortDesc("created_at").Limit(pageSize).Offset(page*pageSize))
	return users, cnt, err
}

// postgres implementation for Verification link repository
type postgresLinkRepo struct {
	repo   r.Repository
	logger *zap.SugaredLogger
}

// Create new repository to work with verification links
func NewVerificationRepository(repo r.Repository) service.VerificationRepository {
	return postgresLinkRepo{repo: repo, logger: zap.S()}
}

//Insert new link
func (rp postgresLinkRepo) Create(ctx context.Context, vl *service.VerificationLink) (*service.VerificationLink, error) {
	err := rp.repo.Insert(ctx, vl)
	return vl, err
}

func (rp postgresLinkRepo) SearchByID(ctx context.Context, id uuid.UUID) (*service.VerificationLink, error) {
	var vl service.VerificationLink
	err := rp.repo.Find(ctx, &vl, r.Where(r.Eq("id", id)))
	return &vl, err
}

//Sets expiration date for a link to current time
func (rp postgresLinkRepo) DeactivateLink(ctx context.Context, id uuid.UUID) error {
	var link *service.VerificationLink
	err := rp.repo.Find(ctx, link, r.Select().Where(r.Eq("id", id)))
	if err != nil {
		return err
	}
	if link == nil {
		return errors.New("could not find the verification link by code")
	}
	link.ExpiresAt = time.Now()
	err = rp.repo.Update(ctx, link)
	return err
}
