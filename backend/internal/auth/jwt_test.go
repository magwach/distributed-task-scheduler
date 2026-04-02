package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	token, err := GenerateToken("user-id", "test@test.com", "admin")

	assert.NoError(t, err)

	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {

	os.Setenv("JWT_SECRET", "test_secret")

	token, err := GenerateToken("user-id", "test@test.com", "admin")

	claims, err := ValidateToken(token)

	assert.NoError(t, err)

	assert.Equal(t, claims.UserID, "user-id")

	assert.Equal(t, claims.Email, "test@test.com")

	assert.Equal(t, claims.Role, "admin")
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("jkbsakeicdqvyucdew")
	assert.Error(t, err)
}
