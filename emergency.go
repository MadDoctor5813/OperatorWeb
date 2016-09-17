package main

import (
    // "log"
    "time"
    
    "gopkg.in/mgo.v2/bson"
)

type Emergency struct {
	Id string `json:"id"`
	Status int `json:"status"`
	Category string `json:"category"`
	Level int `json:"level"`
	InitTime string `json:"initTime"`
	Street string `json:"street"`
	City string `json:"city"`
	Province string `json:"province"`
	PostalCode string `json:"postalCode"`
	Locations Location[] `json:"locations"`
	Notes string `json:"notes"`
	Response string `json:"response"`
	Details string `json:"details"`
	ImageName string `json:"imageName"`
}

type Location struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Street string `json:"street"`
	City string `json:"city"`
	Province string `json:"province"`
	PostalCode string `json:"postalCode"`
}

/*
  ========================================
  Load
  ========================================
*/

func loadEmergencyDB(emergency *Emergency, Id string) error {
    // create new MongoDB session
    collection, session := mongoDBInitialization("emergency")
    defer session.Close()
	
	selector := bson.M{"id" : Id}
    
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

func insertEmergencyDB(emergency *Emergency) (string, error) {
    // create new MongoDB session
    collection, session := mongoDBInitialization("emergency")
    defer session.Close()
    
    newEmergency := new(Emergency)
	newEmergency.Category = emergency.Category
	newEmergency.Details = emergency.Details
	newEmergency.InitTime = time.Now().Format("20060102150405")
	newEmergency.Id = [13:len(bson.NewObjectId().String()) - 2]
	newEmergency.Status = 1
	newEmergency.Level = 0

    // insert resume
    err := collection.Insert(emergency)
    logErrorMessage(err)
    
    return newEmergency.Id, err
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
	
	selector := bson.M{"id" : emergency.Id}
	change := bson.M{"status": emergency.Status, "level": emergency.Level, "notes": emergency.Notes, "response" : emergency.Response}
    update := bson.M{"$set": &change}

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
    selector := bson.M{"id", id}
    err := collection.Remove(selector)
    
    return err
}