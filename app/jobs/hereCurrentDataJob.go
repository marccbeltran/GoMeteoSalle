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

type hereStruct struct {
	Observations struct {
		Location []struct {
			Observation []struct {
				Daylight          string  `json:"daylight"`
				Description       string  `json:"description"`
				SkyInfo           string  `json:"skyInfo"`
				SkyDescription    string  `json:"skyDescription"`
				Temperature       string  `json:"temperature"`
				TemperatureDesc   string  `json:"temperatureDesc"`
				Comfort           string  `json:"comfort"`
				HighTemperature   string  `json:"highTemperature"`
				LowTemperature    string  `json:"lowTemperature"`
				Humidity          string  `json:"humidity"`
				DewPoint          string  `json:"dewPoint"`
				Precipitation1H   string  `json:"precipitation1H"`
				Precipitation3H   string  `json:"precipitation3H"`
				Precipitation6H   string  `json:"precipitation6H"`
				Precipitation12H  string  `json:"precipitation12H"`
				Precipitation24H  string  `json:"precipitation24H"`
				PrecipitationDesc string  `json:"precipitationDesc"`
				AirInfo           string  `json:"airInfo"`
				AirDescription    string  `json:"airDescription"`
				WindSpeed         string  `json:"windSpeed"`
				WindDirection     string  `json:"windDirection"`
				WindDesc          string  `json:"windDesc"`
				WindDescShort     string  `json:"windDescShort"`
				BarometerPressure string  `json:"barometerPressure"`
				BarometerTrend    string  `json:"barometerTrend"`
				Visibility        string  `json:"visibility"`
				SnowCover         string  `json:"snowCover"`
				Icon              string  `json:"icon"`
				IconName          string  `json:"iconName"`
				IconLink          string  `json:"iconLink"`
				AgeMinutes        string  `json:"ageMinutes"`
				ActiveAlerts      string  `json:"activeAlerts"`
				Country           string  `json:"country"`
				State             string  `json:"state"`
				City              string  `json:"city"`
				Latitude          float64 `json:"latitude"`
				Longitude         float64 `json:"longitude"`
				Distance          float64 `json:"distance"`
				Elevation         float64 `json:"elevation"`
			} `json:"observation"`
			Country   string  `json:"country"`
			State     string  `json:"state"`
			City      string  `json:"city"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Distance  float64 `json:"distance"`
			Timezone  int     `json:"timezone"`
		} `json:"location"`
	} `json:"observations"`

	Metric bool `json:"metric"`
}

type HereDataConnection struct{}

func (c HereDataConnection) Run() {

	timeStamp := time.Now().Unix()
	var arrayOfCities = arrayOfCities()
	var totalArray = len(arrayOfCities)

	for i := 0; i < totalArray; i++ {

		latitudeString := strconv.FormatFloat(arrayOfCities[i].Latitude, 'f', 6, 64)
		longitudeString := strconv.FormatFloat(arrayOfCities[i].Longitude, 'f', 6, 64)

		url := "https://weather.cit.api.here.com/weather/1.0/report.json?product=observation&latitude=" + latitudeString + "&longitude=" + longitudeString + "&oneobservation=true&app_id=DemoAppId01082013GAL&app_code=AJKnXv84fjrb0KIHawS0Tg"

		fmt.Print(url)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		response, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(response.Body)

		hereStruct := &hereStruct{}
		json.Unmarshal(body, hereStruct)

		if err != nil {
			log.Fatal(err)
		}

		temperatureString := hereStruct.Observations.Location[0].Observation[0].Temperature
		humidityString := hereStruct.Observations.Location[0].Observation[0].Humidity
		pressureString := hereStruct.Observations.Location[0].Observation[0].BarometerPressure

		idApi := 2

		temperature, err := strconv.ParseFloat(temperatureString, 64)
		humidity, err := strconv.ParseFloat(humidityString, 64)
		pressure, err := strconv.ParseFloat(pressureString, 64)

		stationId := arrayOfCities[i].Station

		station := models.Station{ timeStamp,stationId, arrayOfCities[i].Location, arrayOfCities[i].State,arrayOfCities[i].Latitude, arrayOfCities[i].Longitude, idApi, temperature, humidity, pressure}

		fmt.Print(station,"\n")

		database.Stations.Insert(station)






	}


	database.Stations.RemoveAll(bson.M{"timeStamp": bson.M{"$lt": timeStamp}, "idApi" : 2})



}


func init() {

	revel.OnAppStart(func() {

	//		jobs.Now( HereDataConnection{})
	//		jobs.Schedule("@every 6h", HereDataConnection{})

	})
}


