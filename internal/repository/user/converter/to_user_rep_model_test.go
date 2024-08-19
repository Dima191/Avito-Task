package userrepositoryconverter

import (
	"avito/internal/model"
	"testing"
)

func BenchmarkToUserRepModel(b *testing.B) {
	b.ReportAllocs()

	user := model.User{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToUserRepModel(user)
	}
}
