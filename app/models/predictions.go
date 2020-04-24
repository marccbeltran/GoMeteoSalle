package models

type (
	Prediction struct {
		TimeStamp      int  		`json:"timeStamp" bson:"timeStamp"`
		StationId      int     		`json:"stationId" bson:"stationId"`
		Location       string  		`json:"location" bson:"location"`
		State          int  		`json:"state" bson:"state"`
		Latitude       float64 		`json:"lat"bson:"lat"`
		Longitude      float64 		`json:"long" bson:"long"`
		IdApi          int     		`json:"idApi" bson:"idApi"`
		TemperatureMax float64 		`json:"temperatureMax" bson:"temperatureMax"`
		TemperatureMin float64 		`json:"temperatureMin" bson:"temperatureMin"`
		Humidity       float64 		`json:"humidity" bson:"humidity"`
		Pressure       float64 		`json:"pressure" bson:"pressure"`
	}
)