package househandler

import "fmt"

var (
	APIUrl         = "/api/v1"
	HouseUrl       = "/house"
	CreateHouseUrl = fmt.Sprintf("%s/create", HouseUrl)
)

var (
	LimitQueryParams  = "limit"
	OffsetQueryParams = "offset"
)
