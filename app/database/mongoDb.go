package database

import "gopkg.in/mgo.v2"

/*
Database session
*/
var Session *mgo.Session

/*
Station model connection
*/

var (
	Predictions *mgo.Collection
	Cities      *mgo.Collection
	Stations    *mgo.Collection
)

/*
Init database
*/
func Init(uri, dbname string) error {
	session, err := mgo.Dial(uri)
	if err != nil {
		return err
	}

	//session.SetMode(mgo.Monotonic, true)

	Session = session
	Stations = session.DB(dbname).C("stations")
	Predictions = session.DB(dbname).C("predictions")
	Cities = session.DB(dbname).C("cities")
	return nil
}
