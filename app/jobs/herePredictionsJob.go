package jobs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/marccbeltran/tfmMeteoSalle/app/database"
	"github.com/marccbeltran/tfmMeteoSalle/app/models"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type herePredictionStruct struct {
	DailyForecasts struct {
		ForecastLocation struct {
			Forecast []struct {
				Daylight                 string `json:"daylight"`
				Description              string `json:"description"`
				SkyInfo                  string `json:"skyInfo"`
				SkyDescription           string `json:"skyDescription"`
				TemperatureDesc          string `json:"temperatureDesc"`
				Comfort                  string `json:"comfort"`
				HighTemperature          string `json:"highTemperature"`
				LowTemperature           string `json:"lowTemperature"`
				Humidity                 string `json:"humidity"`
				DewPoint                 string `json:"dewPoint"`
				PrecipitationProbability string `json:"precipitationProbability"`
				PrecipitationDesc        string `json:"precipitationDesc"`
				RainFall                 string `json:"rainFall"`
				SnowFall                 string `json:"snowFall"`
				AirInfo                  string `json:"airInfo"`
				AirDescription           string `json:"airDescription"`
				WindSpeed                string `json:"windSpeed"`
				WindDirection            string `json:"windDirection"`
				WindDesc                 string `json:"windDesc"`
				WindDescShort            string `json:"windDescShort"`
				UvIndex                  string `json:"uvIndex"`
				UvDesc                   string `json:"uvDesc"`
				BarometerPressure        string `json:"barometerPressure"`
				Icon                     string `json:"icon"`
				IconName                 string `json:"iconName"`
				IconLink                 string `json:"iconLink"`
				DayOfWeek                string `json:"dayOfWeek"`
				Weekday                  string `json:"weekday"`
				UtcTime                  string `json:"utcTime"`
			} `json:"forecast"`
			Country   string  `json:"country"`
			State     string  `json:"state"`
			City      string  `json:"city"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Distance  int     `json:"distance"`
			Timezone  int     `json:"timezone"`
		} `json:"forecastLocation"`
	} `json:"dailyForecasts"`
	FeedCreation string `json:"feedCreation"`
	Metric       bool   `json:"metric"`
}


type HereDataPredictionsConnection struct{}

func (c HereDataPredictionsConnection) Run() {

	timeStampDelete := time.Now().Unix()

	arrayOfCities := arrayOfCities()
	totalArray := len(arrayOfCities)

	for j := 0; j < totalArray; j++ {


		latitudeString := strconv.FormatFloat(arrayOfCities[j].Latitude, 'f', 6, 64)
		longitudeString := strconv.FormatFloat(arrayOfCities[j].Longitude, 'f', 6, 64)

		url := "https://weather.cit.api.here.com/weather/1.0/report.json?product=forecast_7days_simple&latitude=" + latitudeString + "&longitude=" + longitudeString + "&oneobservation=true&app_id=DemoAppId01082013GAL&app_code=AJKnXv84fjrb0KIHawS0Tg"

		fmt.Print(url)


		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		response, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(response.Body)

		herePredictionStruct := &herePredictionStruct{}
		json.Unmarshal(body, herePredictionStruct)

		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < 5; i++ {

			utcTime := herePredictionStruct.DailyForecasts.ForecastLocation.Forecast[i].UtcTime
			temperatureMaxString := herePredictionStruct.DailyForecasts.ForecastLocation.Forecast[i].HighTemperature
			temperatureMinString := herePredictionStruct.DailyForecasts.ForecastLocation.Forecast[i].LowTemperature
			pressureString := herePredictionStruct.DailyForecasts.ForecastLocation.Forecast[i].BarometerPressure
			humidityString := herePredictionStruct.DailyForecasts.ForecastLocation.Forecast[i].Humidity
			idApi := 2

			temperatureMax, err := strconv.ParseFloat(temperatureMaxString, 64)

			if err != nil {
				fmt.Print(err.Error())
			}

			date:=utcTime[0:10] + "T12:00:00"

			timeStamp := int(utcTimeToTimeStamp(date))

			temperatureMin, err := strconv.ParseFloat(temperatureMinString, 64)

			if err != nil {
				log.Fatal(err)
			}

			humidity, err := strconv.ParseFloat(humidityString, 64)

			if err != nil {
				log.Fatal(err)
			}

			pressure, err := strconv.ParseFloat(pressureString, 64)

			if err != nil {
				log.Fatal(err)
			}

			prediction := models.Prediction{StationId: arrayOfCities[j].Station,
				Location: arrayOfCities[j].Location,
				State: arrayOfCities[j].State,
				Latitude: arrayOfCities[j].Latitude,
				Longitude: arrayOfCities[j].Longitude,
				IdApi: idApi,
				TimeStamp: timeStamp,
				TemperatureMax: temperatureMax,
				TemperatureMin: temperatureMin,
				Humidity: humidity,
				Pressure: pressure,
			}
			fmt.Print(prediction,"\n")

			database.Predictions.Insert(prediction)
		}


	}


	database.Predictions.RemoveAll(bson.M{"timeStamp": bson.M{"$lt": timeStampDelete }, "idApi" : 2})



}


func utcTimeToTimeStamp(utcTime string) int64 {

	layout:= "2006-01-02T15:04:05"

	t, err := time.Parse(layout, utcTime)
	if err != nil {
		fmt.Println(err)
	}
	timeStamp := (t.Unix())

	return timeStamp

}


func init() {
	revel.OnAppStart(func() {

		//jobs.Now( HereDataPredictionsConnection{})
		jobs.Schedule("@every 48h", HereDataPredictionsConnection{})
	})
}
