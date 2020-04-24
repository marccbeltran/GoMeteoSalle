package models

type (
	Station struct {

		TimeStamp   int64  	 	`json:"timeStamp" bson:"timeStamp"`
		StationId   int     	`json:"stationId" bson:"stationId"`
		Location    string  	`json:"location" bson:"location"`
		State       int  		`json:"state" bson:"state"`
		Latitude    float64 	`json:"lat"bson:"lat"`
		Longitude   float64 	`json:"long" bson:"long"`
		IdApi       int     	`json:"idApi" bson:"idApi"`
		Temperature float64 	`json:"temperature" bson:"temperature"`
		Humidity    float64 	`json:"humidity" bson:"humidity"`
		Pressure    float64 	`json:"pressure" bson:"pressure"`
	}
)
