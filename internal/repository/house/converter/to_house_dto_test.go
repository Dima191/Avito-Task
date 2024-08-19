package houserepositoryconverter

import (
	houserepositorymodel "avito/internal/repository/house/model"
	"testing"
)

func BenchmarkToHouseDto(b *testing.B) {
	b.ReportAllocs()

	house := houserepositorymodel.House{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToHouseDto(house)
	}
}
