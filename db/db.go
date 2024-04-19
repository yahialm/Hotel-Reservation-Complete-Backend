package db

const (
	DBNAME     = "hotel-reservation"
	DBURI      = "mongodb://localhost:27017"
	TestDBNAME = "hotel-reservation-test"
)

type Store struct {
	Hotel   HotelStore
	Room    RoomStore
	User    UserStore
	Booking BookingStore
}
