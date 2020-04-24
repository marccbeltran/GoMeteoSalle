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


type ApixuPredictionStruct struct {
	Alert   struct{} `json:"alert"`
	Current struct {
		Cloud     int `json:"cloud"`
		Condition struct {
			Code int    `json:"code"`
			Icon string `json:"icon"`
			Text string `json:"text"`
		} `json:"condition"`
		FeelslikeC       float64 `json:"feelslike_c"`
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
		TempC            int     `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		Uv               int     `json:"uv"`
		VisKm            int     `json:"vis_km"`
		VisMiles         int     `json:"vis_miles"`
		WindDegree       int     `json:"wind_degree"`
		WindDir          string  `json:"wind_dir"`
		WindKph          float64 `json:"wind_kph"`
		WindMph          float64 `json:"wind_mph"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Astro struct {
				Moonrise string `json:"moonrise"`
				Moonset  string `json:"moonset"`
				Sunrise  string `json:"sunrise"`
				Sunset   string `json:"sunset"`
			} `json:"astro"`
			Date      string `json:"date"`
			DateEpoch int    `json:"date_epoch"`
			Day       struct {
				Avghumidity float64 `json:"avghumidity"`
				AvgtempC    float64 `json:"avgtemp_c"`
				AvgtempF    float64 `json:"avgtemp_f"`
				AvgvisKm    int     `json:"avgvis_km"`
				AvgvisMiles int     `json:"avgvis_miles"`
				Condition   struct {
					Code int    `json:"code"`
					Icon string `json:"icon"`
					Text string `json:"text"`
				} `json:"condition"`
				MaxtempC      float64 `json:"maxtemp_c"`
				MaxtempF      int     `json:"maxtemp_f"`
				MaxwindKph    float64 `json:"maxwind_kph"`
				MaxwindMph    float64 `json:"maxwind_mph"`
				MintempC      float64 `json:"mintemp_c"`
				MintempF      float64 `json:"mintemp_f"`
				TotalprecipIn float64 `json:"totalprecip_in"`
				TotalprecipMm float64 `json:"totalprecip_mm"`
				Uv            float64 `json:"uv"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
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


type apixuDataPredictionsConnectionJob struct{}


func (o apixuDataPredictionsConnectionJob) Run() {

	var timeStampDelete = time.Now().Unix()
	arrayOfCities := arrayOfCities()
	totalArray := len(arrayOfCities)

	for i := 0; i < totalArray; i++ {


		latitudeString := strconv.FormatFloat(arrayOfCities[i].Latitude, 'f', 6, 64)
		longitudeString := strconv.FormatFloat(arrayOfCities[i].Longitude, 'f', 6, 64)
		apikey := "a9de4c35334c4372a96101215191908"

		url := "https://api.apixu.com/v1/forecast.json?key=" + apikey + "&q=" + latitudeString + "," + longitudeString + "&days=5"

		fmt.Print(url)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		response, err := http.Get(url)

		if err != nil {
			fmt.Print(err.Error())

		}

		body, err := ioutil.ReadAll(response.Body)


		ApixuPredictionStruct := &ApixuPredictionStruct{}
		json.Unmarshal(body, ApixuPredictionStruct)

		if err != nil {
			log.Fatal(err)
		}


		for j := 0; j < 4 ; j++ {

			timeStampInt := ApixuPredictionStruct.Forecast.Forecastday[j].Date
			temperatureMax := ApixuPredictionStruct.Forecast.Forecastday[j].Day.MaxtempC
			temperatureMin := ApixuPredictionStruct.Forecast.Forecastday[j].Day.MintempC
			humidity := ApixuPredictionStruct.Forecast.Forecastday[j].Day.Avghumidity
			pressure := ApixuPredictionStruct.Current.PressureMb

			var idApi = 1

			date:= timeStampInt + "T12:00:00"

			timeStamp := int(utcTimeToTimeStamp(date))


			prediction := models.Prediction{StationId: arrayOfCities[i].Station,
											Location: arrayOfCities[i].Location,
											State: arrayOfCities[i].State,
											Latitude: arrayOfCities[i].Latitude,
											Longitude: arrayOfCities[i].Longitude,
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


	database.Predictions.RemoveAll(bson.M{"timeStamp": bson.M{"$lt": timeStampDelete }, "idApi" : 1})




}

func arrayOfCities() []models.City{


	var arrayOfCities []models.City

	pipe:= []bson.M{
		{"$group": bson.M{"_id": "$station",
			"station" : bson.M{"$first": "$station"},
			"location": bson.M{"$first": "$location"},
			"state": bson.M{"$first": "$state"},
			"idApi": bson.M{"$first": "external"},
			"latitude": bson.M{"$first": "$latitude"},
			"longitude": bson.M{"$first": "$longitude"},}}}

	if err := database.Cities.Pipe(pipe).All(&arrayOfCities)

	err != nil {

		log.Fatal(err)
	}

	return arrayOfCities
}

func init() {
	revel.OnAppStart(func() {

		//jobs.Now( apixuDataPredictionsConnectionJob{})
		jobs.Schedule("@every 48h", apixuDataPredictionsConnectionJob{})
	})
}



