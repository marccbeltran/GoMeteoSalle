package controllers

import (
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"log"

)

type Stations struct {
	*revel.Controller
}

func (c Stations) Index() revel.Result {


	pipe:= []bson.M{
			{"$group": bson.M{"_id": "$stationId",
				"stationId": bson.M{"$first": "$stationId"},
				"location": bson.M{"$first": "$location"},
				"state": bson.M{"$first": "$state"},
				"idApi": bson.M{"$first": "external"},
				"latitude": bson.M{"$first": "$lat"},
				"longitude": bson.M{"$first": "$long"},
				"humidity": bson.M{"$avg": "$humidity"},
				"temperature": bson.M{"$avg": "$temperature"},
				"pressure": bson.M{"$avg": "$pressure"},}},{
			"$sort": bson.M{"stationId": -1},
		},}

	var resp []models.StationResponse


	if err := database.Stations.Pipe(pipe).All(&resp)

	err != nil {
		log.Fatal(err)
	}

	return c.RenderJSON(resp)

}
