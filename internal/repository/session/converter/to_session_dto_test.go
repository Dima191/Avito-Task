package sessionrepositoryconverter

import (
	sessionrepositorymodel "avito/internal/repository/session/model"
	"testing"
)

func BenchmarkToSessionDTO(b *testing.B) {
	b.ReportAllocs()

	session := sessionrepositorymodel.Session{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToSessionDTO(session)
	}
}
