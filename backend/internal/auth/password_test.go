package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"

	hash, err := HashPassword(password)

	assert.NoError(t, err)

	assert.NotEmpty(t, hash)

	assert.NotEqual(t, password, hash)
}
