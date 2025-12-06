package infrastructure

import (
	"task9/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHasher_Hash(t *testing.T) {
	hasher := infrastructure.NewBcryptHasher()

	password := "testpassword123"
	hashed, err := hasher.Hash(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestBcryptHasher_Compare(t *testing.T) {
	hasher := infrastructure.NewBcryptHasher()

	password := "testpassword123"
	hashed, err := hasher.Hash(password)
	assert.NoError(t, err)

	t.Run("correct password", func(t *testing.T) {
		result := hasher.Compare(hashed, password)
		assert.True(t, result)
	})

	t.Run("incorrect password", func(t *testing.T) {
		result := hasher.Compare(hashed, "wrongpassword")
		assert.False(t, result)
	})

	t.Run("empty password", func(t *testing.T) {
		result := hasher.Compare(hashed, "")
		assert.False(t, result)
	})
}

func TestBcryptHasher_HashDifferentPasswords(t *testing.T) {
	hasher := infrastructure.NewBcryptHasher()

	password1 := "password1"
	password2 := "password2"

	hashed1, err1 := hasher.Hash(password1)
	hashed2, err2 := hasher.Hash(password2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hashed1, hashed2)
}

func TestBcryptHasher_HashSamePasswordDifferentHashes(t *testing.T) {
	hasher := infrastructure.NewBcryptHasher()

	password := "samepassword"
	hashed1, err1 := hasher.Hash(password)
	hashed2, err2 := hasher.Hash(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hashed1, hashed2)
}

