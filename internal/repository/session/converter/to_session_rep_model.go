package sessionrepositoryconverter

import (
	"avito/internal/model"
	sessionrepositorymodel "avito/internal/repository/session/model"
)

func ToSessionRepModel(session model.Session) sessionrepositorymodel.Session {
	return sessionrepositorymodel.Session{
		SessionID:        session.SessionID,
		UserID:           session.UserID,
		HashRefreshToken: session.HashRefreshToken,
		ExpiresAt:        session.ExpiresAt,
	}
}
