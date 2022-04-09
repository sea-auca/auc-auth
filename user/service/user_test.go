package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sea-auca/auc-auth/user/service"
	"github.com/stretchr/testify/assert"
)

func TestUUIDConversion(t *testing.T) {
	initial := uuid.New()
	str := service.UUIDFromBytes(initial)
	assert.Equal(t, initial, str.Parsed())
}
