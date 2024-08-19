package househandlerconverter

import (
	househandlermodel "avito/internal/handler/house/model"
	"testing"
)

func BenchmarkToHouseDTO(b *testing.B) {
	b.ReportAllocs()

	house := househandlermodel.House{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToHouseDTO(house)
	}
}
