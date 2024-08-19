package sessionrepositoryconverter

import (
	"avito/internal/model"
	sessionrepositorymodel "avito/internal/repository/session/model"
)

func ToSessionDTO(session sessionrepositorymodel.Session) model.Session {
	return model.Session{
		SessionID:        session.SessionID,
		UserID:           session.UserID,
		HashRefreshToken: session.HashRefreshToken,
		ExpiresAt:        session.ExpiresAt,
	}
}
