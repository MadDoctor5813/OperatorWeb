package main

import (
    "encoding/json"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "path"
    "strconv"
    "time"
    "strings"
    
    "github.com/gorilla/context"
    "github.com/gorilla/sessions"
    "github.com/husobee/vestigo"
    "gopkg.in/mgo.v2/bson"
)

type Page struct {
    Name string
}

type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

const (
    errorStatusCode = 398
    serverName = "GWS"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
    router := vestigo.NewRouter()

    // set up router global CORS policy
    router.SetGlobalCors(&vestigo.CorsAccessControl{
        AllowOrigin:      []string{"*"},
        AllowCredentials: false,
        MaxAge:           3600 * time.Second,
    })

    // router.Get("/sandbox/*", staticFile)
    
    fileServerAssets := http.FileServer(http.Dir("assets"))
    router.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Vary", "Accept-Encoding")
        w.Header().Set("Cache-Control", "public, max-age=86400")
        w.Header().Set("Server", serverName)
        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/assets")
        fileServerAssets.ServeHTTP(w, r)
    })
    
    fileServerImages := http.FileServer(http.Dir("images"))
    router.Get("/images/*", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Vary", "Accept-Encoding")
        w.Header().Set("Cache-Control", "public, max-age=86400")
        w.Header().Set("Server", serverName)
        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/images")
        fileServerImages.ServeHTTP(w, r)
    })
    
    fileServerSandbox := http.FileServer(http.Dir("sandbox"))
    router.Get("/sandbox/*", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Vary", "Accept-Encoding")
        w.Header().Set("Cache-Control", "public, max-age=86400")
        w.Header().Set("Server", serverName)
        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/sandbox")
        fileServerSandbox.ServeHTTP(w, r)
    })
    
    // resume
    router.Get("/viewEmergencyJSON/:emergencyId", viewEmergencyJSON)
    
    // login & logout
    router.Post("/checkLoginInfoJSON", checkLoginInfoJSON)
    router.Get("/logoutJSON", logoutJSON)
       
    // resume
    router.Get("/loadEmergencyJSON/:emergencyId", loadEmergencyJSON)
    router.Post("/insertEmergencyJSON/:emergencyId", insertEmergencyJSON)
    router.Post("/updateEmergencyJSON/:emergencyId", updateEmergencyJSON)
    router.Delete("/deleteEmergencyId/:emergencyId", deleteEmergencyId)
 
    // user
    router.Post("/insertUserJSON", insertUserJSON)
    router.Post("/updateUserJSON", updateUserJSON)
    router.Delete("/deleteUserJSON", deleteUserJSON)
    
    // view
    router.Get("/login", viewLogin)
    router.Get("/emergency/:emergencyId", viewEmergency)
    router.Get("/", viewIndex)
   
    log.Println("Listening...")
    if err := http.ListenAndServe(":4242", context.ClearHandler(router)); err != nil {
        log.Println(err)
    }
}

/*
  ========================================
  View
  ========================================
*/

func viewIndex(w http.ResponseWriter, r *http.Request) {
    returnCode := 0

    setHeader(w)
    var homepage Page // placeholder, not used right now
    
    layout := path.Join("templates", "index.html")
    content := path.Join("templates", "content.html")
    
    tmpl, err := template.ParseFiles(layout, content)
    if err != nil {
        returnCode = 1
    }

    if returnCode == 0 {
        if err := tmpl.ExecuteTemplate(w, "my-template", homepage); err != nil {
            returnCode = 2
        }
    }
    
    // error handling
    if returnCode != 0 {
        handleError(returnCode, errorStatusCode, "Index page could not be loaded at this time.", w)
    }
}

func viewEmergency(w http.ResponseWriter, r *http.Request) {
    returnCode := 0

    setHeader(w)
    var homepage Page // placeholder, not used right now
    
    layout := path.Join("templates", "resume.html")
    content := path.Join("templates", "content.html")
    
    tmpl, err := template.ParseFiles(layout, content)
    if err != nil {
        returnCode = 1
    }

    if returnCode == 0 {
        if err := tmpl.ExecuteTemplate(w, "my-template", homepage); err != nil {
            returnCode = 2
        }
    }
    
    // error handling
    if returnCode != 0 {
        handleError(returnCode, errorStatusCode, "Resume page could not be loaded at this time.", w)
    }
}

func viewLogin(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    setHeader(w)
    var homepage Page // placeholder, not used right now
    
    layout := path.Join("templates", "login.html")
    content := path.Join("templates", "content.html")
    
    tmpl, err := template.ParseFiles(layout, content)
    if err != nil {
        returnCode = 1
    }

    if returnCode == 0 {
        if err := tmpl.ExecuteTemplate(w, "my-template", homepage); err != nil {
            returnCode = 2
        }
    }

    // error handling
    if returnCode != 0 {
        handleError(returnCode, errorStatusCode, "Login page could not be loaded at this time.", w)
    }
}

func viewUser(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    if uID, err := readSession("userID", w, r); err == nil && uID != nil {
        setHeader(w)
        var homepage Page // placeholder, not used right now
        var tmpl *template.Template

        layout := path.Join("templates", "user.html")
        content := path.Join("templates", "content.html")

        if tmpl, err = template.ParseFiles(layout, content); err != nil {
            returnCode = 1
        }

        if returnCode == 0 {
            if err = tmpl.ExecuteTemplate(w, "my-template", homepage); err != nil {
                returnCode = 2
            }
        }
        
        // error handling
        if returnCode != 0 {
            handleError(returnCode, errorStatusCode, "User page could not be loaded at this time.", w)
        }
    } else {
        http.Redirect(w, r, "/login", 302)
    }
}

/*
  ========================================
  Resume
  ========================================
*/

func viewEmergencyJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    emergency := new(Emergency)    
    id := vestigo.Param(r, "emergencyId")

    selector := bson.M{"id": id}

    if err := loadEmergencyDB(emergency, id); err != nil {
        returnCode = 1
    }

    if returnCode == 0 {
        if err := json.NewEncoder(w).Encode(emergency); err != nil {
            returnCode = 2
        }
    }

    // error handling
    if returnCode != 0 {
        handleError(returnCode, errorStatusCode, "Emergency could not be loaded at this time.", w)
    }
}

/*
  ========================================
  Login & Logout
  ========================================
*/

func checkLoginInfoJSON(w http.ResponseWriter, r *http.Request) {
    var err error
    returnCode := 0

    user := new(User)
    userID := ""
    
    if err = json.NewDecoder(r.Body).Decode(user); err != nil {
        returnCode = 1
    }
    
    if returnCode == 0 {
        if userID, err = checkPasswordDB(user.Email, user.Password); err != nil {
            returnCode = 2
        }
    }
    
    if returnCode == 0 {
        if userID != "-1" { // if password is correct, then create session
            if err = writeSession("userID", userID, w, r); err != nil {
                returnCode = 3
            }
        }
    }
    
    if returnCode == 0 {
        if err = json.NewEncoder(w).Encode(userID); err != nil {
            returnCode = 4
        }
    }
    
    // error handling
    if returnCode != 0 {
        handleError(returnCode, errorStatusCode, "Incorrect email.", w)
    }
}

func logoutJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    if err := deleteSession(w, r); err != nil {
        returnCode = 1
    }

    if returnCode == 0 {
        if err := json.NewEncoder(w).Encode("logout"); err != nil {
            returnCode = 2
        }
    }
    
    // error handling
    if returnCode != 0 {
        handleError(returnCode, errorStatusCode, "Logout could not be completed at this time.", w)
    }
}

func readSession(key string, w http.ResponseWriter, r *http.Request) (interface{}, error) {
    session, err := store.Get(r, "user-session")
    
    session.Options.MaxAge = 3600 // one hour
    err = session.Save(r, w)
    
    return session.Values[key], err
}

func writeSession(key string, value interface{}, w http.ResponseWriter, r *http.Request) error {
    session, err := store.Get(r, "user-session")
    
    session.Options.MaxAge = 3600 // one hour
    session.Values[key] = value
    err = session.Save(r, w)
    
    return err
}

func deleteSession(w http.ResponseWriter, r *http.Request) error {
    session, err := store.Get(r, "user-session")
    
    session.Options.MaxAge = -1 // delete now
    err = session.Save(r, w)
    
    return err
}

/*
  ========================================
  Resume
  ========================================
*/

func loadEmergencyJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    if uID, err := readSession("userID", w, r); err == nil && uID != nil {
        emergency := new(Emergency)    
        id := vestigo.Param(r, "emergencyId")

        if err = loadEmergencyDB(emergency, id); err != nil {
            returnCode = 1
        }

        if returnCode == 0 {
            if err = json.NewEncoder(w).Encode(emergency); err != nil {
                returnCode = 2
            }
        }

        // error handling
        if returnCode != 0 {
            handleError(returnCode, errorStatusCode, "Emergency could not be loaded at this time.", w)
        }
    } else {
        handleError(3, 403, "Session expired. Please sign in again.", w)
    }
}

func insertEmergencyJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    if uID, err := readSession("userID", w, r); err == nil && uID != nil {
        emergency := new(Emergency)
		json.NewDecoder(r.Body).Decode(emergency)

        if returnCode == 0 {
            if emergency.Id, err = insertEmergencyDB(emergency); err != nil { // Step 1
                returnCode = 1
            }
        }

        if returnCode == 0 {
            if err = json.NewEncoder(w).Encode(emergency); err != nil {
                returnCode = 2
            }
        }

        // error handling
        if returnCode != 0 {
            handleError(returnCode, errorStatusCode, "Emergency could not be inserted at this time.", w)
        }
    } else {
        handleError(3, 403, "Session expired. Please sign in again.", w)
    }
}

func updateEmergencyJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0
    
    if uID, err := readSession("userID", w, r); err == nil && uID != nil {
        emergency := new(Emergency)
		emergency.Id = vestigo.Param(rm "emergencyId")
		json.NewDecoder(r.Body).Decode(emergency)

		if err = updateEmergencyDB(emergency); err != nil {
			returnCode = 1
		}

        if returnCode == 0 {
            if err = json.NewEncoder(w).Encode(emergency); err != nil {
                returnCode = 2
            }
        }

        // error handling
        if returnCode != 0 {
            handleError(returnCode, errorStatusCode, "Emergency could not be updated at this time.", w)
        }
    } else {
        handleError(3, 403, "Session expired. Please sign in again.", w)
    }
}

func deleteEmergencyJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0

    if uID, err := readSession("userID", w, r); err == nil && uID != nil {
        id := vestigo.Param(r, "emergencyId")

        if err := deleteEmergencyDB(id); err != nil {
            returnCode = 1
        }

        // error handling
        if returnCode != 0 {
            handleError(returnCode, errorStatusCode, "Emergency could not be deleted at this time.", w)
        }
    } else {
        handleError(3, 403, "Session expired. Please sign in again.", w)
    }
}

/*
  ========================================
  User
  ========================================
*/

func insertUserJSON(w http.ResponseWriter, r *http.Request) {
    // **** EDIT ****
    
    userID, err := insertUserDB()
    
    err = json.NewEncoder(w).Encode(userID)
    logErrorMessage(err)
}

func updateUserJSON(w http.ResponseWriter, r *http.Request) {
    returnCode := 0

    if uID, err := readSession("userID", w, r); err == nil && uID != nil {
        user := new(User)
        user.UserID = uID.(string)

        if err := json.NewDecoder(r.Body).Decode(user); err != nil {
            returnCode = 1
        }

        if returnCode == 0 {
            if err := updateUSettings(user); err != nil {
                returnCode = 2
            }
        }

        if returnCode == 0 {
            if err := updateRContactTypeName(user.UserID, user.FirstName, user.LastName); err != nil {
                returnCode = 3
            }
        }

        if returnCode == 0 {
            if err := json.NewEncoder(w).Encode(user); err != nil {
                returnCode = 4
            }
        }

        // error handling
        if returnCode != 0 {
            handleError(returnCode, errorStatusCode, "User could not be updated at this time.", w)
        }
    } else {
        handleError(3, 403, "Session expired. Please sign in again.", w)
    }
}

func deleteUserJSON(w http.ResponseWriter, r *http.Request) {
    
    // **** EDIT ****
    
}

/*
  ========================================
  Error
  ========================================
*/

func logErrorMessage(err error) {
    if err != nil {
		log.Println(err)
	}
}

func handleError(returnCode, statusCode int, message string, w http.ResponseWriter) {
    error := new(Error)
    error.Code = returnCode
    error.Message = message

    w.WriteHeader(statusCode)
    if err := json.NewEncoder(w).Encode(error); err != nil {
        log.Println(err)
    }
}

/*
  ========================================
  Basic Functions
  ========================================
*/

func setHeader(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Header().Set("Cache-control", "no-cache, no-store, max-age=0, must-revalidate")
    w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
    w.Header().Set("Server", serverName)
}

func staticFile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Vary", "Accept-Encoding")
    w.Header().Set("Cache-Control", "public, max-age=86400")
    w.Header().Set("Server", serverName)
	http.ServeFile(w, r, r.URL.Path[1:])
}