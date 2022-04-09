package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type VerificationLink struct {
	ID          uuid.UUID `db:",primary"`
	UserID      uuid.UUID
	WasUtilised bool
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (l VerificationLink) Table() string {
	return "users.validation_requests"
}

//Creates new verification link with set parameters
func NewVerificationLink(user uuid.UUID, ttl time.Duration) *VerificationLink {
	vl := &VerificationLink{
		ID:        uuid.New(),
		UserID:    user,
		ExpiresAt: time.Now().Add(ttl),
	}
	return vl
}

type VerificationRepository interface {
	Create(ctx context.Context, vl *VerificationLink) (*VerificationLink, error)
	SearchByID(ctx context.Context, id uuid.UUID) (*VerificationLink, error)
	DeactivateLink(ctx context.Context, id uuid.UUID) error
}
