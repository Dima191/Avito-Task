package apartmentrepositorymodel

type Apartment struct {
	ID               uint32
	ApartmentNumber  int
	HouseID          uint32
	Price            uint32
	NumberOfRooms    uint32
	ModerationStatus string
}
