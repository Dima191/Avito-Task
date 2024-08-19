package apartmentrepositoryconverter

import (
	"avito/internal/model"
	"testing"
)

func BenchmarkToApartmentRepModel(b *testing.B) {
	b.ReportAllocs()

	apartment := model.Apartment{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToApartmentRepModel(apartment)
	}
}
