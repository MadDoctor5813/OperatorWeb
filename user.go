package main

import (
    // "log"

    "gopkg.in/mgo.v2/bson"
)

type User struct {
    UserID        string           `json:"userID"`
    Email         string           `json:"email"`
    Password      string           `json:"password"`
    FirstName     string           `json:"firstName"`
    LastName      string           `json:"lastName"`
}

/*
  ========================================
  Load
  ========================================
*/

func loadSettingsDB(user *User) error {
    // create new MongoDB session
    collection, session := mongoDBInitialization("user")
    defer session.Close()
    
    // retrieve one document with userID as selector
    selector := bson.M{"userid": user.UserID}
    projection := bson.M{"resumeinfo": 0, "portfolioinfo": 0}
    
    err := collection.Find(selector).Select(projection).One(user)
    
    return err
}

/*
  ========================================
  Insert
  ========================================
*/

func insertUserDB() (string, error) {
    // create new MongoDB session
    collection, session := mongoDBInitialization("user")
    defer session.Close()
    
    user := new(User)
    
    // initialize fields // **** TEMPORARY HARD CODE ****
    user.UserID = bson.NewObjectId().String() // get new ObjectId string
    user.UserID = user.UserID[13:len(user.UserID) - 2] // remove extra characters
    user.Email = "ashleycatliu@gmail.com"
    user.Password = "Aa163163"
    user.FirstName = "Ashley"
    user.LastName = "Liu"
    
    // insert user
    err := collection.Insert(user)
    
    return user.UserID, err
}