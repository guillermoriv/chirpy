package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTFlow(t *testing.T) {
	secret := "superSecretKey123"
	wrongSecret := "wrongSecretKey456"
	userID := uuid.New()

	tests := []struct {
		name         string
		userID       uuid.UUID
		signSecret   string
		verifySecret string
		expiresIn    time.Duration
		tamperToken  bool
		wantErr      bool
	}{
		{
			name:         "Valid token roundtrip",
			userID:       userID,
			signSecret:   secret,
			verifySecret: secret,
			expiresIn:    5 * time.Minute,
			tamperToken:  false,
			wantErr:      false,
		},
		{
			name:         "Invalid signature (wrong secret on validation)",
			userID:       userID,
			signSecret:   secret,
			verifySecret: wrongSecret,
			expiresIn:    5 * time.Minute,
			tamperToken:  false,
			wantErr:      true,
		},
		{
			name:         "Expired token",
			userID:       userID,
			signSecret:   secret,
			verifySecret: secret,
			expiresIn:    -1 * time.Minute, // already expired
			tamperToken:  false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(tt.userID, tt.signSecret, tt.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT() unexpected error: %v", err)
			}

			parsedID, err := ValidateJWT(tokenString, tt.verifySecret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && parsedID != tt.userID {
				t.Errorf("ValidateJWT() got userID = %v, want %v", parsedID, tt.userID)
			}
		})
	}
}

func TestValidateJWTMalformedToken(t *testing.T) {
	_, err := ValidateJWT("not.a.valid.jwt.token", "someSecret")
	if err == nil {
		t.Error("ValidateJWT() expected error for malformed token, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		token   string
		wantErr bool
	}{
		{
			name:    "Valid authorization",
			headers: returnHeadersWithAuthorization("Bearer token-1"),
			token:   "token-1",
			wantErr: false,
		},
		{
			name:    "Invalid bearer (wrong formatting)",
			headers: returnHeadersWithAuthorization("Bearertoken-2"),
			token:   "",
			wantErr: true,
		},
		{
			name:    "Invalid bearer (not a token provided)",
			headers: returnHeadersWithAuthorization("Bearer "),
			token:   "",
			wantErr: true,
		},
		{
			name:    "Missing Authorization header",
			headers: http.Header{},
			token:   "",
			wantErr: true,
		},
		{
			name:    "Basic auth header",
			headers: returnHeadersWithAuthorization("Basic dXNlcjpwYXNz"),
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && token != tt.token {
				t.Errorf("GetBearerToken() got token = %v, want %v", token, tt.token)
			}
		})
	}
}

func returnHeadersWithAuthorization(authorization string) http.Header {
	headers := http.Header{}
	headers.Set("Authorization", authorization)
	return headers
}
