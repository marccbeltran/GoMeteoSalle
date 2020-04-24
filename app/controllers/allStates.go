package controllers

import (
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"log"
)



type States struct {
	*revel.Controller
}

func (c States) Index() revel.Result {


	pipe:= []bson.M{
				{"$group": bson.M{"_id": "$state",
					"state": bson.M{"$first": "$state"},
					"idApi": bson.M{"$first": "external"},
					"humidity": bson.M{"$avg": "$humidity"},
					"temperature": bson.M{"$avg": "$temperature"},
					"pressure": bson.M{"$avg": "$pressure"},}}}


	var resp []models.StateResponse

	err := database.Stations.Pipe(pipe).Iter().All(&resp)
	if err != nil {
		log.Fatal(err)
	}

	return c.RenderJSON(resp)

}



