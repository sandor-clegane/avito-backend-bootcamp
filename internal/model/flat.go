package model

// Квартира
type Flat struct {
	ID      int64      `json:"id"`
	HouseID int64      `json:"house_id"`
	Price   int64      `json:"price"`
	Rooms   int64      `json:"rooms"`
	Status  FlatStatus `json:"status"`
}
