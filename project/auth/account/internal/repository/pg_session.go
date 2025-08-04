package account_repo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"erp/api/common"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, userID int64, sessionUUID string, ipAddress, userAgent string, activeCompanyID *int64) (*model.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	DeleteSession(ctx context.Context, token string) error
	DeleteUserSessions(ctx context.Context, userID int64) error
	CleanExpiredSessions(ctx context.Context) error
}

type sessionRepository struct {
	Q  *query.Query
	DB *gorm.DB
	jwtService     helpers.JwtHelper
}

func NewSessionRepository(conn db.Connection,
	helpers *helpers.Helpers) SessionRepository {
	return &sessionRepository{
		Q:  conn.GetQ(),
		DB: conn.GetDB(),
		jwtService: helpers.Jwt,
	}
}

// CreateSession creates a new session for a user
func (r *sessionRepository) CreateSession(ctx context.Context, userID int64,sessionUUID string, ipAddress, userAgent string, activeCompanyID *int64) (*model.Session, error) {
	token, err := r.jwtService.GenerateToken(common.Claims{
		ID:   userID,
		Uuid: sessionUUID,
	})
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		ExpiresAt:       time.Now().Add(24 * time.Hour), // 24 hours from now
		Token:           token,
		CreatedAt:       time.Now(),
		UserID:          userID,
		ActiveCompanyID: activeCompanyID,
	}

	// Set optional fields if provided
	if ipAddress != "" {
		session.IPAddress = &ipAddress
	}
	if userAgent != "" {
		session.UserAgent = &userAgent
	}

	err = r.DB.WithContext(ctx).Create(session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSessionByToken retrieves a session by token
func (r *sessionRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	var session model.Session
	err := r.DB.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession deletes a session by token
func (r *sessionRepository) DeleteSession(ctx context.Context, token string) error {
	return r.DB.WithContext(ctx).Where("token = ?", token).Delete(&model.Session{}).Error
}

// DeleteUserSessions deletes all sessions for a user
func (r *sessionRepository) DeleteUserSessions(ctx context.Context, userID int64) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.Session{}).Error
}

// CleanExpiredSessions removes expired sessions
func (r *sessionRepository) CleanExpiredSessions(ctx context.Context) error {
	return r.DB.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}

// generateToken generates a secure random token for the session
func (r *sessionRepository) generateToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}