package sessions_you_dream_about

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type SessionData struct {
	UserID      string
	CreatedAt   time.Time
	LastActivity time.Time
	IP          string
	Data        map[string]interface{}
}

type SessionManager struct {
	store        Store
	sessionTTL   time.Duration
	cookieName   string
	secureCookie bool
}

func NewManager(store Store, sessionTTL time.Duration, cookieName string, secure bool) *SessionManager {
	if store == nil {
		store = NewInMemoryStore()
	}

	go func() {
		for {
			time.Sleep(time.Hour)
			store.Cleanup(sessionTTL)
		}
	}()

	return &SessionManager{
		store:        store,
		sessionTTL:   sessionTTL,
		cookieName:   cookieName,
		secureCookie: secure,
	}
}

func (sm *SessionManager) CreateSession(userID, ip string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := SessionData{
		UserID:      userID,
		CreatedAt:   time.Now(),
		LastActivity: time.Now(),
		IP:          ip,
		Data:        make(map[string]interface{}),
	}

	if err := sm.store.Save(token, session); err != nil {
		return "", err
	}

	return token, nil
}

func (sm *SessionManager) GetSession(token string) (SessionData, bool) {
	session, err := sm.store.Get(token)
	if err != nil {
		return SessionData{}, false
	}
	return session, true
}

func (sm *SessionManager) UpdateSession(token string, session SessionData) error {
	session.LastActivity = time.Now()
	return sm.store.Save(token, session)
}

func (sm *SessionManager) DestroySession(token string) error {
	return sm.store.Delete(token)
}

func (sm *SessionManager) SessionTimeout() time.Duration {
    return sm.sessionTTL
}

func (sm *SessionManager) CookieName() string {
	return sm.cookieName
}

func (sm *SessionManager) SessionTTL() time.Duration {
	return sm.sessionTTL
}

func (sm *SessionManager) SecureCookie() bool {
	return sm.secureCookie
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
