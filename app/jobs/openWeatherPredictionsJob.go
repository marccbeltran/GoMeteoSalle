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

type OpenWeatherPredictionStruct struct {
	City struct {
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country  string `json:"country"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Timezone int    `json:"timezone"`
	} `json:"city"`
	Cnt  int    `json:"cnt"`
	Cod  string `json:"cod"`
	List []struct {
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Dt    int    `json:"dt"`
		DtTxt string `json:"dt_txt"`
		Main  struct {
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			Temp      float64 `json:"temp"`
			TempKf    float64 `json:"temp_kf"`
			TempMax   float64 `json:"temp_max"`
			TempMin   float64 `json:"temp_min"`
		} `json:"main"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
			ID          int    `json:"id"`
			Main        string `json:"main"`
		} `json:"weather"`
		Wind struct {
			Deg   float64 `json:"deg"`
			Speed float64 `json:"speed"`
		} `json:"wind"`
	} `json:"list"`
	Message float64 `json:"message"`
}


type OpenWeatherConnectionJob struct{}



func (o OpenWeatherConnectionJob) Run() {

	timeStampDelete := time.Now().Unix()
	arrayOfCities := arrayOfCities()
	totalArray := len(arrayOfCities)


	for i := 0; i < totalArray; i++ {

		latitudeString := strconv.FormatFloat(arrayOfCities[i].Latitude, 'f', 6, 64)
		longitudeString := strconv.FormatFloat(arrayOfCities[i].Longitude, 'f', 6, 64)

		url := "https://api.openweathermap.org/data/2.5/forecast?lat=" + latitudeString + "&lon=" + longitudeString + "&units=metric&appid=127f3042ee04e9c7712972b04f1eb66b"

		fmt.Print(url)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		response, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(response.Body)

		OpenWeatherPredictionStruct := &OpenWeatherPredictionStruct{}
		json.Unmarshal(body, OpenWeatherPredictionStruct)

		if err != nil {
			log.Fatal(err)
		}



		for j := 0; j < 40 ; j+=8 {


			timeStampStr := OpenWeatherPredictionStruct.List[j].DtTxt
			temperatureMax := OpenWeatherPredictionStruct.List[j].Main.TempMax
			temperatureMin := OpenWeatherPredictionStruct.List[j].Main.TempMin
			humidityInt := OpenWeatherPredictionStruct.List[j].Main.Humidity
			pressure := OpenWeatherPredictionStruct.List[j].Main.Pressure
			idApi := 3

			humidity := float64(humidityInt)

			date:= timeStampStr[0:10] + "T12:00:00"
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


	database.Predictions.RemoveAll(bson.M{"timeStamp": bson.M{"$lt": timeStampDelete }, "idApi" : 3})


}

func init() {
	revel.OnAppStart(func() {

		//jobs.Now( OpenWeatherConnectionJob{})
		jobs.Schedule("@every 48h", OpenWeatherConnectionJob{})
	})
}
