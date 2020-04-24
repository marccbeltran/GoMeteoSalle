package tests

import (
	"github.com/revel/revel/testing"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestIndexPageWorks() {
	t.Get("/v1")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestAllStatesReturnJson() {
	t.Get("/v1/state")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *AppTest) TestAllStationsReturnJson() {
	t.Get("/v1/stations")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *AppTest) TestAllCitiesReturnJson() {
	t.Get("/v1/cities")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *AppTest) TestAllCitiesResponse200() {
	t.Get("/v1/cities")
	t.AssertOk()
	t.AssertStatus(200)
}

func (t *AppTest) TestAllStationsResponse200() {
	t.Get("/v1/stations")
	t.AssertOk()
	t.AssertStatus(200)
}

func (t *AppTest) TestAllStatesResponse200() {
	t.Get("/v1/state")
	t.AssertOk()
	t.AssertStatus(200)
}

func (t *AppTest) TestOnepredictionResponse200() {
	t.Get("/v1/prediction/1657")
	t.AssertOk()
	t.AssertStatus(200)
}

func (t *AppTest) TestOneStationResponse200() {
	t.Get("/v1/station/1657")
	t.AssertOk()
	t.AssertStatus(200)
}

func (t *AppTest) TestOnePredictionReturnJson() {
	t.Get("/v1/prediction/8096")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *AppTest) TestOneStationReturnJson() {
	t.Get("/v1/station/8096")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *AppTest) TestArrayOfCitiesReturnArray() {
	t.Get("/v1/station/8096")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")

}

func (t *AppTest) After() {
	println("Tear down")
}


