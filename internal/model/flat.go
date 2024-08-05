package model

// Квартира
type Flat struct {
	ID            int64      `json:"id"`
	HouseID       int64      `json:"house_id"`
	Price         int64      `json:"price"`
	NumberOfRooms int64      `json:"rooms"`
	Status        FlatStatus `json:"status"`
}
