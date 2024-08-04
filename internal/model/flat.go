package model

// Квартира
type Flat struct {
	ID            int        `json:"id"`
	HouseID       int        `json:"house_id"`
	FlatNumber    string     `json:"apartment_number"`
	Price         int        `json:"price"`
	NumberOfRooms int        `json:"number_of_rooms"`
	Status        FlatStatus `json:"flat_status"`
}
