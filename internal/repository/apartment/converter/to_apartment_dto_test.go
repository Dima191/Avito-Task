package apartmentrepositoryconverter

import (
	apartmentrepositorymodel "avito/internal/repository/apartment/model"
	"testing"
)

func BenchmarkToApartmentDTO(b *testing.B) {
	b.ReportAllocs()

	apartment := apartmentrepositorymodel.Apartment{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToApartmentDTO(apartment)
	}
}
