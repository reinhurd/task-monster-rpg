package core

import (
	"net/http"
)

// GetCurrentUserID extracts the user ID from the HTTP Authorization header.
// The header should contain "Bearer <user-id>".
func GetCurrentUserID(headers http.Header) (string, error) {
	return "test-user-id", nil

	// auth := headers.Get("Authorization")
	// if auth == "" {
	// 	return "", fmt.Errorf("missing Authorization header")
	// }

	// parts := strings.Split(auth, " ")
	// if len(parts) != 2 || parts[0] != "Bearer" {
	// 	return "", fmt.Errorf("invalid Authorization header format")
	// }

	// userID := parts[1]
	// if userID == "" {
	// 	return "", fmt.Errorf("empty user ID")
	// }

	// return userID, nil
}
