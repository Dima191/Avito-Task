package apartmenthandler

import "fmt"

var (
	APIUrl             = "/api/v1"
	ApartmentsUrl      = "/apartment"
	CreateApartmentUrl = fmt.Sprintf("%s/create", ApartmentsUrl)

	ApartmentID        = "apartment_id"
	ApartmentUrl       = fmt.Sprintf("%s/{%s}", ApartmentsUrl, ApartmentID)
	UpdateApartmentUrl = fmt.Sprintf("%s/update", ApartmentUrl)

	HouseID                = "house_id"
	ApartmentsByHouseIDUrl = fmt.Sprintf("%s/{%s}", ApartmentsUrl, HouseID)
)

var (
	LimitQueryParams  = "limit"
	OffsetQueryParams = "offset"
)
