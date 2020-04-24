package models

type (
	StationResponse struct {
		StationId   int     	`json:"stationId" bson:"stationId"`
		Location    string  	`json:"location" bson:"location"`
		State       int 		`json:"state" bson:"state"`
		Latitude    float64		`json:"latitude"bson:"latitude"`
		Longitude   float64 	`json:"longitude" bson:"longitude"`
		IdApi       string  	`json:"idApi" bson:"idApi"`
		Temperature float64 	`json:"temperature" bson:"temperature"`
		Humidity    float64 	`json:"humidity" bson:"humidity"`
		Pressure    float64 	`json:"pressure" bson:"pressure"`
	}
)
