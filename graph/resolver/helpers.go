package resolver

import (
	"context"
	"errors"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

const (
	adminRole = "admin"
)

// GetUserIDFromContext functions to extraact user info from GraphQL context
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return 0, ErrUnauthorized
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, ErrUnauthorized
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role := ctx.Value("user_role")
	if role == nil {
		return "", ErrUnauthorized
	}

	if r, ok := role.(string); ok {
		return r, nil
	}

	return "", ErrUnauthorized
}

func IsAdminFromContext(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == adminRole
}
