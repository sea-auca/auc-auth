package service

import (
	"context"
	"time"

	"github.com/go-rel/rel"
	"github.com/google/uuid"
)

type UUID string

func (id UUID) Parsed() uuid.UUID {
	return uuid.MustParse(string(id))
}

func UUIDFromBytes(id uuid.UUID) UUID {
	return UUID(id.String())
}

//Main user model struct
type User struct {
	Email        string
	Fullname     string
	Hash         string      //bcrypt hash
	Password     string      `db:"-"` // virtual field that will automatically converted to hash
	UUID         UUID        `db:"uuid,primary"`
	AccessLevels Permissions `db:"permissions"`
	Active       bool        //deactivated account are *deleted*
	Verified     bool        //flag is set automatically after verification of email
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const AucaEmail = `([a-z]+)(_){1}[a-z]{1,4}(@auca.kg|@alumni.auca.kg)`

func (u User) Table() string {
	return "user_space.users"
}

func NewUser(email string) *User {
	return &User{Email: email, UUID: UUID(uuid.NewString()), AccessLevels: None, Active: true}
}

type UserRepository interface {
	//cruds

	Create(ctx context.Context, u *User, ch ...rel.Mutator) (*User, error)
	Update(ctx context.Context, u *User) error
	RelUpdate(ctx context.Context, u *User, ch ...rel.Mutator) error // this breaks clean architecture completely!
	//getters

	GetByID(ctx context.Context, id UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	PaginatedView(ctx context.Context, page, pageSize int) ([]*User, int, error)
}
