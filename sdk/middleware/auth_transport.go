package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type BrifleClientState struct {
	Token             string       // The token used to authenticate the client with the Brifle API.
	LastAuthenticated int64        // The timestamp of the last successful authentication.
	AuthInterval      int64        // The interval in seconds between automatic re-authentication attempts.
	Credentials       *Credentials // Optional credentials for the client, if needed.
}

type Credentials struct {
	ApiKey    string // Username for authentication
	ApiSecret string // Password for authentication
}

type AuthTransport struct {
	BaseTransport     http.RoundTripper
	State             BrifleClientState
	tokenMutex        sync.RWMutex
	AllowTokenRenewal bool // Flag to allow or disallow token renewal
	RenewToken        func() (string, error)
}

const LOGIN_PATH = "/v1/auth/login"

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Use base transport (default if nil)
	transport := t.BaseTransport
	if transport == nil {
		transport = http.DefaultTransport
	}

	if req.URL.Path == LOGIN_PATH {
		// If the request is for the login path, do not use the token
		return transport.RoundTrip(req)
	}

	currentTime := time.Now().UTC().UnixMilli() // Assuming current time is passed in context
	// Lock the mutex to safely access and modify the state
	t.tokenMutex.RLock()
	needUpdate := t.State.Token == "" ||
		t.State.LastAuthenticated == 0 ||
		t.State.AuthInterval == 0 ||
		currentTime-t.State.LastAuthenticated >= t.State.AuthInterval*1000
	t.tokenMutex.RUnlock()

	// Check if the token needs to be renewed
	if needUpdate && t.AllowTokenRenewal {
		// Lock the mutex for writing to update the token
		t.tokenMutex.Lock()
		newToken, err := t.RenewToken()
		if err != nil {
			return nil, fmt.Errorf("token renewal failed: %w", err)
		}
		t.State.Token = newToken
		t.State.LastAuthenticated = currentTime
		t.tokenMutex.Unlock()
	}

	// Safely read the token
	t.tokenMutex.RLock()
	token := t.State.Token
	t.tokenMutex.RUnlock()

	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Set("Authorization", "Bearer "+token)

	return transport.RoundTrip(clonedReq)
}
