package models

type (
	City struct {
		Station   		int    		 	`json:"station" bson:"station"`
		Location  		string 			`json:"location" bson:"location"`
		State     		int  	 		`json:"state" bson:"state"`
		Latitude  		float64 		`json:"latitude"bson:"latitude"`
		Longitude 		float64 		`json:"longitude" bson:"longitude"`
		postalCode 		int     		`json:"postalCode" bson:"postalCode"`
	}
)