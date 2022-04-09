package service

import (
	"context"
	"regexp"
	"time"

	"github.com/google/uuid"
)

//auca email regex
const aucaEmail = `([a-z]+)(_){1}[a-z]{1,4}(@auca.kg|@alumni.auca.kg)`

//Main user model struct
type User struct {
	ID          uuid.UUID `db:",primary"`
	Email       string
	IsActive    bool //deactivated account are *deleted*
	IsValidated bool //flag is set automatically after verification of email
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserData struct {
	UserID          uuid.UUID `db:"user_id,primary"`
	FirstName       string
	LastName        string
	AvatarURL       string
	Department      string
	YearOfAdmission int
	BirthDate       time.Time
	UpdatedAt       time.Time
}

type AuthSettings struct {
	UserID                uuid.UUID `db:"user_id,primary"`
	RefreshTokenRetention int       `db:"refresh_token_retention_in_hours"`
	TokenRetention        int       `db:"main_token_retention_in_minutes"`
	EnforceTwoFactorAuth  bool
	UpdatedAt             time.Time
}

func (u User) Table() string {
	return `users.users`
}

func NewUser(email string) *User {
	return &User{Email: email, ID: uuid.New()}
}

type UserRepository interface {
	//cruds

	Create(ctx context.Context, u *User) (*User, error)
	Update(ctx context.Context, u *User) error
	//getters

	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	PaginatedView(ctx context.Context, page, pageSize int) ([]*User, int, error)
}

func IsAucaEmail(email string) bool {
	match, _ := regexp.MatchString(aucaEmail, email)
	return match
}
