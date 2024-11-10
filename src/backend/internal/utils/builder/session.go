package builder

import (
	"time"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type SessionBuilder struct {
	clientID  uuid.UUID
	sessionID uuid.UUID
	ttl       time.Time
}

func NewSessionBuilder() *SessionBuilder {
	return &SessionBuilder{}
}

func (b *SessionBuilder) SetClientID(clientID uuid.UUID) *SessionBuilder {
	b.clientID = clientID
	return b
}

func (b *SessionBuilder) SetSessionID(sessionID uuid.UUID) *SessionBuilder {
	b.sessionID = sessionID
	return b
}

func (b *SessionBuilder) SetTTL(ttl time.Time) *SessionBuilder {
	b.ttl = ttl
	return b
}

func (b *SessionBuilder) Build() entity.Session {
	return entity.Session{
		ClientID:  b.clientID,
		SessionID: b.sessionID,
		TTL:       b.ttl,
	}
}
