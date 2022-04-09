package service

import (
	"context"
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
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (l VerificationLink) Table() string {
	return "user_space.verify_links"
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

type VerificationRepository interface {
	Create(ctx context.Context, vl *VerificationLink) (*VerificationLink, error)
	SearchByCode(ctx context.Context, code string) (*VerificationLink, error)
	SearchByUser(ctx context.Context, id UUID) ([]*VerificationLink, error)
	DeactivateLink(ctx context.Context, uuid UUID, code string) error
	DeactivateAllLinks(ctx context.Context, uuid UUID) error
}
