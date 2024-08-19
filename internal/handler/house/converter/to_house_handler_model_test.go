package househandlerconverter

import (
	"avito/internal/model"
	"testing"
)

func BenchmarkToHouseHandlerModel(b *testing.B) {
	b.ReportAllocs()

	house := model.House{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToHouseHandlerModel(house)
	}
}
