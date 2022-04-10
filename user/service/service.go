package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNonAucaEmail   = errors.New("attemt to register with non-auca email")
	ErrNotImplemented = errors.New("unimplemented function")
)

type UserService interface {
	RegisterUser(ctx context.Context, email string) error
	DeactivateAccount(ctx context.Context, u *User) error
	ReactivateAccount(ctx context.Context, email string) error
	ValidateUser(ctx context.Context, validatonID uuid.UUID) error

	//Read operations

	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*User, int, error)
}

type service struct {
	repo     UserRepository
	lg       *zap.SugaredLogger
	linkRepo VerificationRepository
}

func NewService(repo UserRepository, vRepo VerificationRepository) UserService {
	return service{
		repo:     repo,
		lg:       zap.S().With("service", "user"),
		linkRepo: vRepo,
	}
}

// Creates a placeholder record for user and send an email for verification
func (s service) RegisterUser(ctx context.Context, email string) error {
	if !IsAucaEmail(email) {
		s.lg.Info("Non-Auca email registration", zap.String("email", email))
		return ErrNonAucaEmail
	}
	var newUser = NewUser(email)
	newUser, err := s.repo.Create(ctx, newUser)
	if err != nil {
		s.lg.Errorw("failed to create user", "error", err, "user", newUser)
		return err
	}
	vl := NewVerificationLink(newUser.ID, time.Hour*24*7)
	vl, err = s.linkRepo.Create(ctx, vl)
	if err != nil {
		s.lg.Errorw("failed to create invite link", "error", err)
		return err
	}
	err = s.sendVereficationEmail(ctx, vl.ID, email)
	if err != nil {
		s.lg.Errorw("failed to send the email for invite", "error", err)
	}
	return nil
}

// Suspends user's account. It can be later reactivated
func (s service) DeactivateAccount(ctx context.Context, u *User) error {
	u.IsActive = false
	err := s.repo.Update(ctx, u)
	if err != nil {
		s.lg.Errorw("failed to update the activity status", "error", err)
	}
	return err
}

// Restores user's account after suspension
func (s service) ReactivateAccount(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		s.lg.Errorw("could not find user with such email", "error", err, "email", email)
		return err
	}

	if user.IsActive {
		err = errors.New("user is already active")
		s.lg.Errorw("failed to reactivate user", "error", err)
		return err
	}

	link := NewVerificationLink(user.ID, time.Hour*6)
	link, err = s.linkRepo.Create(ctx, link)
	if err != nil {
		s.lg.Errorw("failed to create reactivation link", "error", err)
		return err
	}
	//err = s.sendReactivationEmail(link.Link, email)
	if err != nil {
		s.lg.Errorw("failed to send the email for invite", "error", err)
	}
	return err
}

// Returns the user with specified uuid
func (s service) GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error) {
	user, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		s.lg.Errorw("failed to fetch user with uuid", zap.String("id", ID.String()), zap.Error(err))
	}
	return user, err
}

// Returns the user with specified email
func (s service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.lg.Error("failed to fetch user with email", zap.String("email", email), zap.Error(err))
	}
	return user, err
}

// Returns the pagination data
func (s service) ListUsers(ctx context.Context, page, pageSize int) ([]*User, int, error) {
	return s.repo.PaginatedView(ctx, page, pageSize)
}

// Checks the verification code and activates the user account
func (s service) ValidateUser(ctx context.Context, validationID uuid.UUID) error {
	vl, err := s.linkRepo.SearchByID(ctx, validationID)
	if err != nil {
		s.lg.Error("failed to verify user, link is abscent", zap.Error(err))
		return err
	}
	if vl.WasUtilised || time.Now().After(vl.ExpiresAt) {
		return errors.New("link is invalid")
	}
	user, err := s.repo.GetByID(ctx, vl.UserID)
	if err != nil {
		return err
	}
	user.IsActive = true
	user.IsValidated = true
	err = s.repo.Update(ctx, user)
	err = s.linkRepo.DeactivateLink(ctx, vl.ID)
	if err != nil {
		s.lg.Errorw("failed to verify user", "reason", err, "user", user.ID)
	}
	return err
}
