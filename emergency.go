package main

import (
	// "log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Emergency struct {
	Id          string     `json:"id"`
	Status      int        `json:"status"`
	Category    string     `json:"category"`
	Level       int        `json:"level"`
	InitTime    string     `json:"initTime"`
	Street      string     `json:"street"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	PostalCode  string     `json:"postalCode"`
	Locations   []Location `json:"locations"`
	Notes       string     `json:"notes"`
	Response    string     `json:"response"`
	Details     string     `json:"details"`
	Description string     `json:"description"`
	ImageName   string     `json:"imageName"`
}

type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Street     string  `json:"street"`
	City       string  `json:"city"`
	Province   string  `json:"province"`
	PostalCode string  `json:"postalCode"`
	Time       string  `json:"time"`
}

/*
  ========================================
  Load
  ========================================
*/

func loadEmergenciesDB(emergencies *[]Emergency, status int) error {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	selector := bson.M{"status": status}

	// retrieve one document with resumeID as selector
	err := collection.Find(selector).All(emergencies)
	logErrorMessage(err)

	return err
}

func loadEmergencyDB(emergency *Emergency, Id string) error {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	selector := bson.M{"id": Id}

	// retrieve one document with resumeID as selector
	err := collection.Find(selector).One(emergency)
	logErrorMessage(err)

	return err
}

/*
  ========================================
  Insert
  ========================================
*/

func testInsertEmergency() {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	emergency := new(Emergency)
	emergency.Category = "Cakesplosion"
	emergency.Details = "I was happy. Then my cake exploded. Now I'm sad."
	emergency.InitTime = time.Now().Format("20060102150405")
	emergency.Id = bson.NewObjectId().String()
	emergency.Id = emergency.Id[13 : len(emergency.Id)-2]
	emergency.Status = 1
	emergency.Level = 0

	// insert resume
	err := collection.Insert(emergency)
	logErrorMessage(err)
}

func insertEmergencyDB(emergency *Emergency) (string, error) {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	emergency.Category = emergency.Category
	emergency.InitTime = time.Now().Format("20060102150405")
	emergency.Id = bson.NewObjectId().String()
	emergency.Id = emergency.Id[13 : len(emergency.Id)-2]
	emergency.Status = 1
	emergency.Level = 0

	// insert resume
	err := collection.Insert(emergency)
	logErrorMessage(err)

	return emergency.Id, err
}

/*
  ========================================
  Update
  ========================================
*/

func updateEmergencyDBUser(emergency *Emergency) error {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	selector := bson.M{"id": emergency.Id}
	change := bson.M{"details": emergency.Details, "description": emergency.Description, "imagename": emergency.ImageName}
	update := bson.M{"$set": &change}

	err := collection.Update(selector, update)

	return err
}

func updateEmergencyDBAdmin(emergency *Emergency) error {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	selector := bson.M{"id": emergency.Id}
	change := bson.M{"status": emergency.Status, "level": emergency.Level, "notes": emergency.Notes, "response": emergency.Response}
	update := bson.M{"$set": &change}

	err := collection.Update(selector, update)

	return err
}

func updateLocationDB(location *Location, emergencyId string) error {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	selector := bson.M{"id": emergencyId}
	change := bson.M{"locations": bson.M{"latitude": location.Latitude, "longitude": location.Longitude, "street": location.Street, "city": location.City, "province": location.Province, "postalcode": location.PostalCode, "time": location.Time}}
	update := bson.M{"$addToSet": &change}

	err := collection.Update(selector, update)

	return err
}

func testInsertLocation() {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	emergencyId := "57dd54a08a46bb871d000001"

	location := new(Location)
	location.Latitude = 43.6489778
	location.Longitude = -79.3799492
	location.Street = "44 King St W"
	location.City = "Toronto"
	location.Province = "ON"
	location.PostalCode = "M5H 1H1"
	location.Time = time.Now().Format("20060102150405")

	selector := bson.M{"id": emergencyId}
	change := bson.M{"locations": bson.M{"latitude": location.Latitude, "longitude": location.Longitude, "street": location.Street, "city": location.City, "province": location.Province, "postalcode": location.PostalCode, "time": location.Time}}
	update := bson.M{"$addToSet": &change}

	// update location
	err := collection.Update(selector, update)
	logErrorMessage(err)

	change = bson.M{"street": location.Street, "city": location.City, "province": location.Province, "postalcode": location.PostalCode}
	update = bson.M{"$set": &change}

	err = collection.Update(selector, update)
	logErrorMessage(err)
}

/*
  ========================================
  Delete
  ========================================
*/

func deleteEmergencyDB(id string) error {
	// create new MongoDB session
	collection, session := mongoDBInitialization("emergency")
	defer session.Close()

	// find document and delete resume
	selector := bson.M{"id": id}
	err := collection.Remove(selector)

	return err
}
