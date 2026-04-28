package utils_test

import (
	"testing"

	"github.com/m0xyu/learning-go-shop/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestPassword_HashAndCheck(t *testing.T) {
	originalPassword := "super_secret_password_123"

	hashedPassword, err := utils.HashPassword(originalPassword)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, originalPassword, hashedPassword)

	isMatch := utils.CheckPassword(originalPassword, hashedPassword)

	assert.True(t, isMatch)

	isMatchWrong := utils.CheckPassword("wrong_password", hashedPassword)

	assert.False(t, isMatchWrong)
}
