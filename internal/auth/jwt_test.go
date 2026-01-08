package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  http.Header
		want    string
		wantErr bool
	}{
		{
			name: "valid bearer token",
			header: http.Header{
				"Authorization": []string{"Bearer my-token-123"},
			},
			want:    "my-token-123",
			wantErr: false,
		},
		{
			name:    "missing authorization header",
			header:  http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing bearer prefix",
			header: http.Header{
				"Authorization": []string{"my-token-123"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty bearer token",
			header: http.Header{
				"Authorization": []string{"Bearer "},
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "valid JWT creation",
			userID:      uuid.New(),
			tokenSecret: "my-secret-key",
			expiresIn:   time.Hour,
			wantErr:     false,
		},
		{
			name:        "valid JWT with short expiration",
			userID:      uuid.New(),
			tokenSecret: "another-secret",
			expiresIn:   time.Minute,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("MakeJWT() returned empty token")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	tests := []struct {
		name        string
		setupToken  func() string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := MakeJWT(userID, secret, time.Hour)
				return token
			},
			tokenSecret: secret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name: "wrong secret",
			setupToken: func() string {
				token, _ := MakeJWT(userID, secret, time.Hour)
				return token
			},
			tokenSecret: "wrong-secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name: "expired token",
			setupToken: func() string {
				// Create a token that expires immediately
				token, _ := MakeJWT(userID, secret, -time.Hour)
				return token
			},
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name: "invalid token format",
			setupToken: func() string {
				return "not-a-valid-jwt-token"
			},
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			gotUserID, err := ValidateJWT(token, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestMakeAndValidateJWT_Integration(t *testing.T) {
	t.Run("create and validate token successfully", func(t *testing.T) {
		userID := uuid.New()
		secret := "integration-test-secret"
		expiresIn := time.Hour

		// Create a token
		token, err := MakeJWT(userID, secret, expiresIn)
		if err != nil {
			t.Fatalf("MakeJWT() failed: %v", err)
		}

		// Validate the token
		validatedUserID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("ValidateJWT() failed: %v", err)
		}

		// Check that the user ID matches
		if validatedUserID != userID {
			t.Errorf("ValidateJWT() returned wrong userID: got %v, want %v", validatedUserID, userID)
		}
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		userID := uuid.New()
		secret := "expired-test-secret"

		// Create a token that's already expired
		token, err := MakeJWT(userID, secret, -time.Second)
		if err != nil {
			t.Fatalf("MakeJWT() failed: %v", err)
		}

		// Try to validate the expired token
		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("ValidateJWT() should have rejected expired token")
		}
	})

	t.Run("token with wrong secret is rejected", func(t *testing.T) {
		userID := uuid.New()
		correctSecret := "correct-secret"
		wrongSecret := "wrong-secret"

		// Create a token with the correct secret
		token, err := MakeJWT(userID, correctSecret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() failed: %v", err)
		}

		// Try to validate with wrong secret
		_, err = ValidateJWT(token, wrongSecret)
		if err == nil {
			t.Error("ValidateJWT() should have rejected token signed with different secret")
		}
	})

	t.Run("multiple users have different tokens", func(t *testing.T) {
		user1 := uuid.New()
		user2 := uuid.New()
		secret := "multi-user-secret"

		// Create tokens for both users
		token1, err := MakeJWT(user1, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() failed for user1: %v", err)
		}

		token2, err := MakeJWT(user2, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() failed for user2: %v", err)
		}

		// Tokens should be different
		if token1 == token2 {
			t.Error("Different users should have different tokens")
		}

		// Validate both tokens
		validated1, err := ValidateJWT(token1, secret)
		if err != nil {
			t.Fatalf("ValidateJWT() failed for token1: %v", err)
		}
		if validated1 != user1 {
			t.Errorf("Token1 validated to wrong user: got %v, want %v", validated1, user1)
		}

		validated2, err := ValidateJWT(token2, secret)
		if err != nil {
			t.Fatalf("ValidateJWT() failed for token2: %v", err)
		}
		if validated2 != user2 {
			t.Errorf("Token2 validated to wrong user: got %v, want %v", validated2, user2)
		}
	})
}
