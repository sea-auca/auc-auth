package user

import (
	"context"
	"errors"
	"time"

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

func ChangeUser(user interface{}, p params.Params) *changeset.Changeset {
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
	changeset.ValidateRequired(ch, []string{"email", "verified"})
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

// postgres implementation for Verification link repository
type postgresLinkRepo struct {
	repo   r.Repository
	logger *zap.SugaredLogger
}

// Create new repository to work with verification links
func NewVerificationRepository(repo r.Repository) VerificationRepository {
	return postgresLinkRepo{repo: repo, logger: zap.S()}
}

//Insert new link
func (rp postgresLinkRepo) Create(ctx context.Context, vl *VerificationLink) (*VerificationLink, error) {
	err := rp.repo.Insert(ctx, vl)
	return vl, err
}

func (rp postgresLinkRepo) SearchByCode(ctx context.Context, code string) (*VerificationLink, error) {
	var vl *VerificationLink
	err := rp.repo.Find(ctx, vl, r.Where(r.Eq("code", code)))
	return vl, err
}

func (rp postgresLinkRepo) SearchByUser(ctx context.Context, id UUID) ([]*VerificationLink, error) {
	var vl []*VerificationLink
	err := rp.repo.FindAll(ctx, vl, r.Select().SortDesc("created_at").Where(r.Eq("user_id", id)))
	return vl, err
}

//Sets expiration date for a link to current time
func (rp postgresLinkRepo) DeactivateLink(ctx context.Context, uuid UUID, code string) error {
	var link *VerificationLink
	err := rp.repo.Find(ctx, link, r.Select().Where(r.Eq("code", code), r.Eq("user_id", uuid)))
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

//Deactivate all links for specified user
func (rp postgresLinkRepo) DeactivateAllLinks(ctx context.Context, uuid UUID) error {
	vls, err := rp.SearchByUser(ctx, uuid)
	if err != nil {
		return err
	}
	var codes []string
	for _, v := range vls {
		codes = append(codes, v.Link)
	}
	_, err = rp.repo.UpdateAny(ctx,
		r.From("user_space.verify_links").Where(r.In("code", codes)),
		r.Set("expires_at", time.Now()),
	)
	return err
}
