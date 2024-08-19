package apartmenthandlerconverter

import (
	apartmenthandlermodel "avito/internal/handler/apartment/model"
	"testing"
)

func BenchmarkToApartmentDTO(b *testing.B) {
	b.ReportAllocs()

	apartment := apartmenthandlermodel.Apartment{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToApartmentDTO(apartment)
	}
}
