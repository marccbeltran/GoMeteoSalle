package models

type (
	StateResponse struct {
		State       string  	`json:"state" bson:"state"`
		IdApi       string  	`json:"idApi" bson:"idApi"`
		Temperature float64 	`json:"temperature" bson:"temperature"`
		Humidity    float64 	`json:"humidity" bson:"humidity"`
		Pressure    float64 	`json:"pressure" bson:"pressure"`
	}
)
