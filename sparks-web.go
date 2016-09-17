package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/husobee/vestigo"
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
	serverName      = "GWS"
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

	fileServerCSS := http.FileServer(http.Dir("css"))
	router.Get("/css/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", serverName)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/css")
		fileServerCSS.ServeHTTP(w, r)
	})

	fileServerHTML := http.FileServer(http.Dir("html"))
	router.Get("/html/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", serverName)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/html")
		fileServerHTML.ServeHTTP(w, r)
	})

	fileServerJS := http.FileServer(http.Dir("js"))
	router.Get("/js/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", serverName)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/js")
		fileServerJS.ServeHTTP(w, r)
	})

	fileServerPKG := http.FileServer(http.Dir("pkg"))
	router.Get("/pkg/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", serverName)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/pkg")
		fileServerPKG.ServeHTTP(w, r)
	})

	fileServerIMG := http.FileServer(http.Dir("img"))
	router.Get("/img/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", serverName)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/img")
		fileServerIMG.ServeHTTP(w, r)
	})

	// login & logout
	router.Post("/checkLoginInfoJSON", checkLoginInfoJSON)
	router.Get("/logoutJSON", logoutJSON)

	// emergency
	router.Get("/loadEmergenciesJSON/:status", loadEmergenciesJSON)
	router.Get("/loadEmergencyJSON/:emergencyId", loadEmergencyJSON)
	router.Post("/insertEmergencyJSON", insertEmergencyJSON)
	router.Post("/updateEmergencyJSON/:emergencyId", updateEmergencyJSON)
	router.Post("/updateLocationJSON/:emergencyId", updateLocationJSON)
	router.Delete("/deleteEmergencyId/:emergencyId", deleteEmergencyJSON)

	// view
	router.Get("/pending", viewAdmin)
	router.Get("/in-progress", viewAdmin)
	router.Get("/complete", viewAdmin)
	router.Get("/archives", viewAdmin)
	router.Get("/trash", viewAdmin)
	router.Get("/sign-in", viewLogin)
	router.Get("/", viewLogin)

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

func viewAdmin(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	setHeader(w)
	var homepage Page // placeholder, not used right now

	layout := path.Join("html", "admin.html")
	content := path.Join("html", "content.html")

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
		handleError(returnCode, errorStatusCode, "Admin page could not be loaded at this time.", w)
	}
}

func viewLogin(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	setHeader(w)
	var homepage Page // placeholder, not used right now

	layout := path.Join("html", "sign-in.html")
	content := path.Join("html", "content.html")

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
  Emergency
  ========================================
*/

func loadEmergenciesJSON(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	if uID, err := readSession("userID", w, r); err == nil && uID != nil {
		var emergencies []Emergency
		var statusInt int
		var err error

		statusStr := vestigo.Param(r, "status")
		if statusInt, err = strconv.Atoi(statusStr); err != nil {
			returnCode = 1
		}

		if returnCode == 0 {
			if err = loadEmergenciesDB(&emergencies, statusInt); err != nil {
				returnCode = 1
			}
		}

		if returnCode == 0 {
			if err = json.NewEncoder(w).Encode(&emergencies); err != nil {
				returnCode = 2
			}
		}

		// error handling
		if returnCode != 0 {
			handleError(returnCode, errorStatusCode, "Emergencies could not be loaded at this time.", w)
		}
	} else {
		handleError(3, 403, "Session expired. Please sign in again.", w)
	}
}

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

		if err := json.NewDecoder(r.Body).Decode(emergency); err != nil {
			returnCode = 1
		}

		if returnCode == 0 {
			if emergency.Id, err = insertEmergencyDB(emergency); err != nil { // Step 1
				returnCode = 2
			}
		}

		if returnCode == 0 {
			if err = json.NewEncoder(w).Encode(emergency); err != nil {
				returnCode = 3
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
		emergency.Id = vestigo.Param(r, "emergencyId")

		if err := json.NewDecoder(r.Body).Decode(emergency); err != nil {
			returnCode = 1
		}

		if returnCode == 0 {
			if err = updateEmergencyDB(emergency); err != nil {
				returnCode = 2
			}
		}

		if returnCode == 0 {
			if err = json.NewEncoder(w).Encode(emergency); err != nil {
				returnCode = 3
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

func updateLocationJSON(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	if uID, err := readSession("userID", w, r); err == nil && uID != nil {
		location := new(Location)
		emergencyId := vestigo.Param(r, "emergencyId")

		if json.NewDecoder(r.Body).Decode(location); err != nil {
			returnCode = 1
		}

		if returnCode == 0 {
			if err = updateLocationDB(location, emergencyId); err != nil {
				returnCode = 2
			}
		}

		if returnCode == 0 {
			if err = json.NewEncoder(w).Encode(location); err != nil {
				returnCode = 3
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
