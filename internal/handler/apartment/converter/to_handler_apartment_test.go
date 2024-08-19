package apartmenthandlerconverter

import (
	"avito/internal/model"
	"testing"
)

func BenchmarkToHandlerModelApartment(b *testing.B) {
	b.ReportAllocs()

	apartment := model.Apartment{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToHandlerModelApartment(apartment)
	}
}
