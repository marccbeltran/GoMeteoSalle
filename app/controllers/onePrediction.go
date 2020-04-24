package controllers

import (
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)

type Prediction struct {
	*revel.Controller
}

	func (c Prediction) Index() revel.Result {

		stationId2 :=  c.Params.Route.Get("stationid")
		stationInt, _ := strconv.Atoi(stationId2)



		pipe:= []bson.M{{"$match": bson.M{"stationId": stationInt}},
			{"$group": bson.M{"_id": "$timeStamp",
				"stationId": bson.M{"$first": "$stationId"},
				"timeStamp": bson.M{"$first": "$timeStamp"},
				"location": bson.M{"$first": "$location"},
				"state": bson.M{"$first": "$state"},
				"idApi": bson.M{"$first": "external"},
				"latitude": bson.M{"$first": "$lat"},
				"longitude": bson.M{"$first": "$long"},
				"humidity": bson.M{"$avg": "$humidity"},
				"temperature": bson.M{"$avg": "$temperatureMax"},
				"pressure": bson.M{"$avg": "$pressure"},}},{
			"$sort": bson.M{"timeStamp": 1},
		},{
				"$limit": 5,
			}}

		var resp []models.PredictionResponse

		err := database.Predictions.Pipe(pipe).Iter().All(&resp)

		if err != nil {
			log.Fatal(err)
		}

		return c.RenderJSON(resp)

	}



