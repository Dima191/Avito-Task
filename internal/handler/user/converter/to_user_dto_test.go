package userhandlerconverter

import (
	userhandlermodel "avito/internal/handler/user/model"
	"testing"
)

func BenchmarkToUserDto(b *testing.B) {
	b.ReportAllocs()

	user := userhandlermodel.User{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToUserDto(user)
	}
}
