package sessionrepositoryconverter

import (
	"avito/internal/model"
	"testing"
)

func BenchmarkToSessionRepModel(b *testing.B) {
	b.ReportAllocs()

	session := model.Session{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToSessionRepModel(session)
	}
}
