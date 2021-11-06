package user_test

import (
	"sea/auth/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDConversion(t *testing.T) {
	initial := uuid.New()
	str := user.UUIDFromBytes(initial)
	assert.Equal(t, initial, str.Parsed())
}
