package controllers

import (
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Cities struct {
	*revel.Controller
}

func (c Cities ) Index() revel.Result {

	var results []models.City

	if err := database.Cities.Find(bson.M{}).All(&results); err != nil {

		c.Response.Status = http.StatusServiceUnavailable

	}

	return c.RenderJSON(results)


}