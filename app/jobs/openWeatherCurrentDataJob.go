package jobs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type OpenweatherStruct struct {
	Coord struct {
		Lon int `json:"lon"`
		Lat int `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Rain struct {
		ThreeH float64 `json:"3h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

type openWeatherDataConnection struct{}

func (c openWeatherDataConnection) Run() {



	timeStamp := time.Now().Unix()
	arrayOfCities := arrayOfCities()
	totalArray := len(arrayOfCities)
	fmt.Print(totalArray)



	for i := 0; i < totalArray; i++ {


		var latitudeString = fmt.Sprintf("%.5f", arrayOfCities[i].Latitude)
		var longitudeString = fmt.Sprintf("%.5f", arrayOfCities[i].Longitude)

		url := "https://api.openweathermap.org/data/2.5/weather?lat=" + latitudeString + "&lon=" + longitudeString + "&units=metric&appid=127f3042ee04e9c7712972b04f1eb66b"

		fmt.Print(url)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		response, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(response.Body)

		OpenweatherStruct := &OpenweatherStruct{}
		json.Unmarshal(body, OpenweatherStruct)

		if err != nil {
			log.Fatal(err)
		}

		temperature := OpenweatherStruct.Main.Temp
		humidityInt := OpenweatherStruct.Main.Humidity
		pressureInt := OpenweatherStruct.Main.Pressure

		idApi := 3

		humidity := float64(humidityInt)
		pressure := float64(pressureInt)

		stationId := arrayOfCities[i].Station


		station := models.Station{ timeStamp,stationId, arrayOfCities[i].Location, arrayOfCities[i].State,arrayOfCities[i].Latitude, arrayOfCities[i].Longitude, idApi, temperature, humidity, pressure}


		fmt.Print(station,"\n")

		database.Stations.Insert(station)

	}

	database.Stations.RemoveAll(bson.M{"timeStamp": bson.M{"$lt": timeStamp}, "idApi" : 3})



}

func init() {
	revel.OnAppStart(func() {
	//	jobs.Now( openWeatherDataConnection{})
	//	jobs.Schedule("@every 6h", openWeatherDataConnection{})
	})
}