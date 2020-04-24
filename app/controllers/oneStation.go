package controllers

import (
	_ "fmt"
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)



type Station struct {
	*revel.Controller
}

func (c Station) Index() revel.Result {


	stationId2 :=  c.Params.Route.Get("stationid")
	stationInt, _ := strconv.Atoi(stationId2)

	pipe:= []bson.M{{"$match": bson.M{"stationId": stationInt}},
					{"$group": bson.M{"_id": "$stationId",
						"stationId": bson.M{"$first": "$stationId"},
						"location": bson.M{"$first": "$location"},
						"state": bson.M{"$first": "$state"},
						"idApi": bson.M{"$first": "external"},
						"latitude": bson.M{"$first": "$lat"},
						"longitude": bson.M{"$first": "$long"},
						"humidity": bson.M{"$avg": "$humidity"},
						"temperature": bson.M{"$avg": "$temperature"},
						"pressure": bson.M{"$avg": "$pressure"},}}}


	var resp []models.StationResponse

	err := database.Stations.Pipe(pipe).Iter().All(&resp)
	if err != nil {
		log.Fatal(err)
	}

	return c.RenderJSON(resp)

}



