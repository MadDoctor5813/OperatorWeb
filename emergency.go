package main

import (
	// "log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Emergency struct {
	Id         string     `json:"id"`
	Status     int        `json:"status"`
	Category   string     `json:"category"`
	Level      int        `json:"level"`
	InitTime   string     `json:"initTime"`
	Street     string     `json:"street"`
	City       string     `json:"city"`
	Province   string     `json:"province"`
	PostalCode string     `json:"postalCode"`
	Locations  []Location `json:"locations"`
	Notes      string     `json:"notes"`
	Response   string     `json:"response"`
	Details    string     `json:"details"`
	ImageName  string     `json:"imageName"`
}

type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Street     string  `json:"street"`
	City       string  `json:"city"`
	Province   string  `json:"province"`
	PostalCode string  `json:"postalCode"`
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
	emergency.Category = "Pillow Crack"
	emergency.Details = "My pillow cracked and Deloitte Digital won't give use a pillow."
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
	emergency.Details = emergency.Details
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

func updateEmergencyDB(emergency *Emergency) error {
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
	change := bson.M{"locations": bson.M{"latitude": location.Latitude, "longitude": location.Longitude, "street": location.Street, "city": location.City, "province": location.Province, "postalcode": location.PostalCode}}
	update := bson.M{"$addToSet": &change}

	err := collection.Update(selector, update)

	return err
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
