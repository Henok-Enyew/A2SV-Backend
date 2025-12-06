package infrastructure

import (
	"os"
	"task9/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTGenerator_Generate(t *testing.T) {
	generator := infrastructure.NewJWTGenerator()

	userID := "507f1f77bcf86cd799439011"
	username := "testuser"
	role := "admin"

	token, err := generator.Generate(userID, username, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTGenerator_Validate(t *testing.T) {
	generator := infrastructure.NewJWTGenerator()

	userID := "507f1f77bcf86cd799439011"
	username := "testuser"
	role := "admin"

	token, err := generator.Generate(userID, username, role)
	assert.NoError(t, err)

	t.Run("valid token", func(t *testing.T) {
		claims, err := generator.Validate(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims["user_id"])
		assert.Equal(t, username, claims["username"])
		assert.Equal(t, role, claims["role"])
	})

	t.Run("invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		_, err := generator.Validate(invalidToken)
		assert.Error(t, err)
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := generator.Validate("")
		assert.Error(t, err)
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := generator.Validate("not.a.valid.jwt.token")
		assert.Error(t, err)
	})
}

func TestJWTGenerator_DifferentSecrets(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret1")
	generator1 := infrastructure.NewJWTGenerator()

	userID := "507f1f77bcf86cd799439011"
	username := "testuser"
	role := "admin"

	token, err := generator1.Generate(userID, username, role)
	assert.NoError(t, err)

	os.Setenv("JWT_SECRET", "secret2")
	generator2 := infrastructure.NewJWTGenerator()

	_, err = generator2.Validate(token)
	assert.Error(t, err)

	os.Unsetenv("JWT_SECRET")
}

func TestJWTGenerator_ExpiredToken(t *testing.T) {
	generator := infrastructure.NewJWTGenerator()

	userID := "507f1f77bcf86cd799439011"
	username := "testuser"
	role := "admin"

	token, err := generator.Generate(userID, username, role)
	assert.NoError(t, err)

	claims, err := generator.Validate(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
}

