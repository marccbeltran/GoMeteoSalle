# API METEOSALLE

REST API to create and retrieve weather predictions with Go and Revel framework and MongoDB

# Install MongoDB 

To run this API you need a installation of MongoDB Database with this configuration and collections

<code>database.uri  = "mongodb://localhost:27017"</code>

<code>database.name = "meteosalle"</code>

    Predictions *mgo.Collection
	Cities      *mgo.Collection
	Stations    *mgo.Collection



# Install dependencies

This project uses Revel framework and DEP, to install dependencies you can execute this command

<code>DEP ensure</code>

# Run

Once Revel framework is installed, you can run the server by:

<code>revel run</code>

Note that the project must be located under <code>$GOPATH/src/tfmMeteoSalle</code>

# Routes

The API routes are defined in <code>conf/routes</code> file:

    
    MODULE  /@jobs
    MODULE  /@testRunner

    GET     /v1/                                      
    GET     /v1/stations/                              
    GET     /v1/cities/                             
    GET     /v1/states/                               
    GET     /v1/prediction/:stationid                  
    GET     /v1/station/:stationid                     