package main

import (
	"context"
	"strings"

	"encore.dev/beta/auth"
)

//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, error) {
	// Dev-friendly: accept any non-empty token
	t := strings.TrimSpace(token)
	if t == "" {
		return "", auth.ErrUnauthenticated
	}
	return auth.UID(strings.ToLower(t)), nil
}

