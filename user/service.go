package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-rel/changeset/params"
	"go.uber.org/zap"
	"gopkg.in/mail.v2"
)

var ErrNotImplemented = errors.New("unimplemented function")

type UserService interface {
	// CUD ops

	RegisterUser(ctx context.Context, email string) (*User, error)
	DeactivateAccount(ctx context.Context, u *User) error
	ReactivateAccount(ctx context.Context, email string) error
	UpdateUser(ctx context.Context, u *User, p params.Params) error

	//Read operations

	GetUserByID(ctx context.Context, ID UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*User, int, error)

	// sensative data update

	VerifyUser(ctx context.Context, link string) error
	RequestNewPassword(ctx context.Context, email string) error
	SetNewPassword(ctx context.Context, newPassword, resetCode string) error
}

type service struct {
	repo     UserRepository
	lg       *zap.SugaredLogger
	email    mail.Sender
	linkRepo VerificationRepository
}

func NewService(repo UserRepository, email mail.Sender, vRepo VerificationRepository) UserService {
	return service{
		repo:     repo,
		lg:       zap.S().With("service", "user"),
		email:    email,
		linkRepo: vRepo,
	}
}

func (s service) sendVereficationEmail(code, email string) error {
	const hi_msg = `Hello! <hr/> Thank you for joining the AU Cloud. We are happy to see you there. `
	const invite = `Please follow this link and complete your registration: <br/>`
	var link = "https://sea.auca.kg/user/verify?code=" + code
	var href = fmt.Sprintf("<a href=\"%s\">%s<a/>", link, link)
	m := mail.NewMessage()
	m.SetHeader("From", "sea@auca.kg")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to AU Cloud!")
	m.SetBody("text/html", hi_msg+invite+href)

	return s.email.Send("sea@auca.kg", []string{"student_s@auca.kg"}, m)
}

// Creates a placeholder record for user and send an email for verification
func (s service) RegisterUser(ctx context.Context, email string) (*User, error) {
	newUser := NewUser(email)
	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		s.lg.Errorw("failed to create user", "error", err)
		return nil, err
	}

	vl := NewVerificationLink(user.UUID, time.Hour*24*7, false)
	vl, err = s.linkRepo.Create(ctx, vl)
	if err != nil {
		s.lg.Errorw("failed to create invite link", "error", err)
		return nil, err
	}
	err = s.sendVereficationEmail(vl.Link, email)
	if err != nil {
		s.lg.Errorw("failed to send the email for invite", "error", err)
	}
	return user, nil
}

// Suspends user's account. It can be later reactivated
func (s service) DeactivateAccount(ctx context.Context, u *User) error {
	return ErrNotImplemented
}

// Restores user's account after suspension
func (s service) ReactivateAccount(ctx context.Context, email string) error {
	return ErrNotImplemented
}

// Updates non-system data on user
func (s service) UpdateUser(ctx context.Context, u *User, p params.Params) error {
	return ErrNotImplemented
}

// Returns the user with specified uuid
func (s service) GetUserByID(ctx context.Context, ID UUID) (*User, error) {
	return nil, ErrNotImplemented
}

// Returns the user with specified email
func (s service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, ErrNotImplemented
}

// Returns the pagination data
func (s service) ListUsers(ctx context.Context, page, pageSize int) ([]*User, int, error) {
	return nil, 0, ErrNotImplemented
}

// Checks the verification code and activates the user account
func (s service) VerifyUser(ctx context.Context, link string) error {
	return ErrNotImplemented
}

// Sends an email with a link for password reset
func (s service) RequestNewPassword(ctx context.Context, email string) error {
	return ErrNotImplemented
}

// sets new password for user after password reset
func (s service) SetNewPassword(ctx context.Context, newPassword, resetCode string) error {
	return ErrNotImplemented
}
