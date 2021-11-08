package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-rel/changeset/params"
	"github.com/go-rel/rel"
	"go.uber.org/zap"
	"gopkg.in/mail.v2"
)

var ErrNotImplemented = errors.New("unimplemented function")

type UserService interface {
	// CUD ops

	RegisterUser(ctx context.Context, email string) error
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
	var link = "https://sea.auca.kg/user/verify?code=" + code + "&action=i"
	var href = fmt.Sprintf("<a href=\"%s\">%s<a/>", link, link)
	m := mail.NewMessage()
	m.SetHeader("From", "sea@auca.kg")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to AU Cloud!")
	m.SetBody("text/html", hi_msg+invite+href)

	return s.email.Send("sea@auca.kg", []string{email}, m)
}

// Creates a placeholder record for user and send an email for verification
func (s service) RegisterUser(ctx context.Context, email string) error {
	newUser := NewUser(email)
	changeset := ChangeUser(newUser, params.Map{})
	if changeset.Error() != nil {
		s.lg.Errorw("failed to create new user", "reason", changeset.Errors())
		return changeset.Error()
	}
	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		s.lg.Errorw("failed to create user", "error", err)
		return err
	}

	vl := NewVerificationLink(user.UUID, time.Hour*24*7, false)
	vl, err = s.linkRepo.Create(ctx, vl)
	if err != nil {
		s.lg.Errorw("failed to create invite link", "error", err)
		return err
	}
	err = s.sendVereficationEmail(vl.Link, email)
	if err != nil {
		s.lg.Errorw("failed to send the email for invite", "error", err)
	}
	return nil
}

// Suspends user's account. It can be later reactivated
func (s service) DeactivateAccount(ctx context.Context, u *User) error {
	changeset := rel.NewChangeset(u)
	u.Active = false
	err := s.repo.RelUpdate(ctx, u, changeset)
	if err != nil {
		s.lg.Errorw("failed to update the activity status", "error", err)
	}
	return err
}

func (s service) sendReactivationEmail(code, email string) error {
	const hi = "Hello! Thank you for reactivating you AU Cloud account. <hr/> Follow this link and your account will be reactivated"
	link := "https://sea.auca.kg/user/verify?code=" + code + "&action=r"
	var link_full = fmt.Sprintf("<br/><a href=\"%s\">%s<a/>", link, link)
	m := mail.NewMessage()
	m.SetHeader("From", "sea@auca.kg")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Account reactivation")
	m.SetBody("text/html", hi+link_full)

	return s.email.Send("sea@auca.kg", []string{email}, m)
}

// Restores user's account after suspension
func (s service) ReactivateAccount(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		s.lg.Errorw("could not find user with such email", "error", err, "email", email)
		return err
	}

	if user.Active {
		err = errors.New("user is already active")
		s.lg.Errorw("failed to reactivate user", "error", err)
		return err
	}

	link := NewVerificationLink(user.UUID, time.Hour*6, false)
	link, err = s.linkRepo.Create(ctx, link)
	if err != nil {
		s.lg.Errorw("failed to create reactivation link", "error", err)
		return err
	}
	err = s.sendReactivationEmail(link.Link, email)
	if err != nil {
		s.lg.Errorw("failed to send the email for invite", "error", err)
	}
	return err
}

// Updates non-system data on user
func (s service) UpdateUser(ctx context.Context, u *User, p params.Params) error {
	changeset := ChangeUser(u, p)
	if changeset.Error() != nil {
		s.lg.Errorw("failed to update user, changeset error", "errors", changeset.Errors())
	}
	err := s.repo.RelUpdate(ctx, u, changeset)
	if err != nil {
		s.lg.Errorw("failed to update user's data", "error", err)
	}
	return err
}

// Returns the user with specified uuid
func (s service) GetUserByID(ctx context.Context, ID UUID) (*User, error) {
	user, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		s.lg.Errorw("failed to fetch user with uuid", "uuid", ID, "error", err)
	}
	return user, err
}

// Returns the user with specified email
func (s service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.lg.Errorw("failed to fetch user with email", "email", email, "error", err)
	}
	return user, err
}

// Returns the pagination data
func (s service) ListUsers(ctx context.Context, page, pageSize int) ([]*User, int, error) {
	return s.repo.PaginatedView(ctx, page, pageSize)
}

// Checks the verification code and activates the user account
func (s service) VerifyUser(ctx context.Context, link string) error {
	vl, err := s.linkRepo.SearchByCode(ctx, link)
	if err != nil {
		s.lg.Errorw("failed to verify user, link is abscent", "reason", err, "link code", link)
		return err
	}
	user, err := s.repo.GetByID(ctx, vl.UserID)
	if err != nil {
		return err
	}
	changeset := rel.NewChangeset(user)
	user.Verified = true
	user.Active = true
	err = s.repo.RelUpdate(ctx, user, changeset)
	if err != nil {
		s.lg.Errorw("failed to verify user", "reason", err, "user", user.UUID)
	}
	return err
}

func (s service) sendPasswordResetLink(code, email string) error {
	const hi = "Hello! <hr/> Follow this link to reset your password"
	link := "https://sea.auca.kg/user/reset?code=" + code + "&action=r"
	var link_full = fmt.Sprintf("<br/><a href=\"%s\">%s<a/>", link, link)
	m := mail.NewMessage()
	m.SetHeader("From", "sea@auca.kg")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password reset")
	m.SetBody("text/html", hi+link_full)

	return s.email.Send("sea@auca.kg", []string{email}, m)
}

// Sends an email with a link for password reset
func (s service) RequestNewPassword(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		s.lg.Errorw("could not find user with such email", "error", err, "email", email)
		return err
	}

	link := NewVerificationLink(user.UUID, time.Hour*12, true)
	link, err = s.linkRepo.Create(ctx, link)
	if err != nil {
		s.lg.Errorw("failed to create password reset link", "error", err)
		return err
	}
	err = s.sendPasswordResetLink(link.Link, email)
	if err != nil {
		s.lg.Errorw("failed to send the email for password reset", "error", err)
	}
	return err
}

// sets new password for user after password reset
func (s service) SetNewPassword(ctx context.Context, newPassword, resetCode string) error {
	vl, err := s.linkRepo.SearchByCode(ctx, resetCode)
	if err != nil {
		s.lg.Errorw("failed to get rest code", "reason", err, "link code", resetCode)
		return err
	}
	user, err := s.repo.GetByID(ctx, vl.UserID)
	if err != nil {
		return err
	}
	changeset := ChangeUser(user, params.Map{"password": newPassword})
	if changeset.Error() != nil {
		s.lg.Errorw("failed to create new password, changeset error", "errors", changeset.Errors())
		return changeset.Error()
	}
	err = s.repo.RelUpdate(ctx, user, changeset)
	if err != nil {
		s.lg.Errorw("failed to verify user", "reason", err, "user", user.UUID)
	}
	return err
}
