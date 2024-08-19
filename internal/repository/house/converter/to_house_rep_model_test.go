package houserepositoryconverter

import (
	"avito/internal/model"
	"testing"
)

func BenchmarkToHouseRepModel(b *testing.B) {
	b.ReportAllocs()

	house := model.House{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToHouseRepModel(house)
	}
}
