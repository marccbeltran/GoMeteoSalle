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
	"strconv"
	"time"
)

type ApixuStruct struct {
	Current struct {
		Cloud     int `json:"cloud"`
		Condition struct {
			Code int    `json:"code"`
			Icon string `json:"icon"`
			Text string `json:"text"`
		} `json:"condition"`
		FeelslikeC       int     `json:"feelslike_c"`
		FeelslikeF       float64 `json:"feelslike_f"`
		GustKph          float64 `json:"gust_kph"`
		GustMph          float64 `json:"gust_mph"`
		Humidity         int     `json:"humidity"`
		IsDay            int     `json:"is_day"`
		LastUpdated      string  `json:"last_updated"`
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		PrecipIn         int     `json:"precip_in"`
		PrecipMm         int     `json:"precip_mm"`
		PressureIn       float64 `json:"pressure_in"`
		PressureMb       float64 `json:"pressure_mb"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		Uv               int     `json:"uv"`
		VisKm            int     `json:"vis_km"`
		VisMiles         int     `json:"vis_miles"`
		WindDegree       int     `json:"wind_degree"`
		WindDir          string  `json:"wind_dir"`
		WindKph          int     `json:"wind_kph"`
		WindMph          float64 `json:"wind_mph"`
	} `json:"current"`
	Location struct {
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Localtime      string  `json:"localtime"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Lon            float64 `json:"lon"`
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		TzID           string  `json:"tz_id"`
	} `json:"location"`
}


type ApixuDataConnection struct{}


func (a ApixuDataConnection) Run() {


	timeStamp := time.Now().Unix()
	arrayOfCities := arrayOfCities()
	totalArray := len(arrayOfCities)


	for i := 0; i < totalArray; i++ {

		latitudeString := strconv.FormatFloat(arrayOfCities[i].Latitude, 'f', 4, 64)
		longitudeString := strconv.FormatFloat(arrayOfCities[i].Longitude, 'f', 4, 64)
		apikey := "a9de4c35334c4372a96101215191908"

		url := "https://api.apixu.com/v1/current.json?key=" + apikey + "&q=" + latitudeString + "," + longitudeString

		fmt.Print(url)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		response, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
		}

		ApixuStruct := &ApixuStruct{}

		json.Unmarshal(body, ApixuStruct)



		temperatureInt := ApixuStruct.Current.TempC
		humidityInt := ApixuStruct.Current.Humidity
		pressureInt := ApixuStruct.Current.PressureMb

		idApi := 1


		temperature := float64(temperatureInt)
		humidity := float64(humidityInt)
		pressure := float64(pressureInt)

		stationId := arrayOfCities[i].Station


		station := models.Station{ timeStamp,stationId, arrayOfCities[i].Location, arrayOfCities[i].State,arrayOfCities[i].Latitude, arrayOfCities[i].Longitude, idApi, temperature, humidity, pressure}

		fmt.Print(station, "\n")


		database.Stations.Insert(station)


	}


	database.Stations.RemoveAll(bson.M{"timeStamp": bson.M{"$lt": timeStamp}, "idApi" : 1})



}

func init() {

	revel.OnAppStart(func() {

	//		jobs.Now( ApixuDataConnection{})
	//		jobs.Schedule("@every 6h", ApixuDataConnection{})
	})
}
