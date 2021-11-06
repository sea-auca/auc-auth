package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type VerificationLink struct {
	ID              int64 `db:",primary"`
	UserID          UUID
	Link            string
	IsPasswordReset bool
	ExpiresAt       time.Time
	CreateAt        time.Time
	UpdatedAt       time.Time
}

//Creates new verification link with set parameters
func NewVerificationLink(user UUID, ttl time.Duration, isPwd bool) *VerificationLink {
	//80 bit random link from uuid
	link := strings.Join(strings.Split(uuid.NewString(), "-")[:5], "")
	vl := &VerificationLink{
		UserID:          user,
		Link:            link,
		IsPasswordReset: isPwd,
		ExpiresAt:       time.Now().Add(ttl),
	}
	return vl
}
