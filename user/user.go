package user

import (
	"time"

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
func (u User) Table() string {
	return "user_space.users"
}
