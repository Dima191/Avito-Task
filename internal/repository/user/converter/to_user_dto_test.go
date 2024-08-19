package userrepositoryconverter

import (
	userrepositorymodel "avito/internal/repository/user/model"
	"testing"
)

func BenchmarkToUserDTO(b *testing.B) {
	b.ReportAllocs()

	user := userrepositorymodel.User{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToUserDTO(user)
	}
}
