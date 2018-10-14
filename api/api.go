package api

import (
	"encoding/json"
	"fmt"
	."jvh_local/IMT2681-assignement-1/data"
	"github.com/marni/goigc"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)
// Unique id for each track
var idCount int = 0
// Array containing each id of track
var ids = make([]int, 0)
// Map of tracks and their id
var tracks = make(map[string]Tracks)
// The time
var t = time.Now()

func conver(d time.Duration) string {
	// For string manipulation
	var felles []string
	sec := d.Seconds()

	const (
		mins	= 60         // Minutes in seconds
		hours   = 3600	   // Hours in seconds
		days    = 86400	   // Days in seconds
		months  = 2629746  // Months in seconds
		years   = 31556952 // Years in seconds
	)


	felles = append(felles, "P")

	// Divide seconds with years in seconds to find number of current years
	year  := int(sec / years)
	if year >= 1 {
		felles = append(felles, strconv.Itoa(year))
		felles = append(felles, "Y")
		// Subtracting the number of years in seconds - to provide right amount of seconds
		sec -= float64(years * year)
	}
	// Divide seconds with months in seconds to find number of current months
	month := int(sec / months)
	if month >= 1 {
		felles = append(felles, strconv.Itoa(month))
		felles = append(felles, "M")
		// Subtracting the number of months in seconds - to provide right amount of seconds
		sec -= float64(months * month)
	}

	// Divide seconds with days in seconds to find number of current days
	day   := int(sec / days)	 // Days in seconds
	if day >= 1 {
		felles = append(felles, strconv.Itoa(day))
		felles = append(felles, "D")
		// Subtracting the number of days in seconds - to provide right amount of seconds
		sec -= float64(86400 * day)
	}

	felles = append(felles, "T")

	// Divide seconds with hours in seconds to find number of current hours
	hour  := int(sec / hours) 	 // Hours in seconds
	if hour >= 1 {
		felles = append(felles, strconv.Itoa(hour))
		felles = append(felles, "H")
		// Subtracting the number of hours in seconds - to provide right amount of seconds
		sec -= float64(hours * hour)

	}

	// Divide seconds with minutes in seconds to find number of current minutes
	min   := int(sec / mins) 		 // Minutes in seconds
	if min >= 1 {
		felles = append(felles, strconv.Itoa(min))
		felles = append(felles, "M")
		sec -= float64(mins * min)

	}

	if sec >= 0 {
		felles = append(felles, strconv.Itoa(int(sec)))
		felles = append(felles, "S")
	}

	// Joins the part of the slice to one string
	k := strings.Join(felles, "")
	// Returns string with corresponding timestamp
	return k
}



func InfoHandler(w http.ResponseWriter, r *http.Request) {

	check := regexp.MustCompile("^/igcinfo/api/$")
	// Checks if the URL match the regex
	if check.MatchString(r.URL.Path) {
		// Time since application started
		uptime := time.Since(t)
		iso := conver(uptime)
		infoApi := Info{
			iso,
			"Service for IGC tracks.",
			"v1",
		}

		// Set the header to json
		w.Header().Set("Content-Type", "application/json")
		// Encodes information to user
		json.NewEncoder(w).Encode(infoApi)
	} else {
		// If unknown path is provided - StatusCode 404 is returned
		http.NotFound(w,r)
		//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// Based on request method
	case "POST":
		PostAPI(w, r)

	case "GET":
		GetAPI(w, r)
	default:
		http.Error(w, "Not implemented yet", http.StatusNotImplemented)
	}

}

func PostAPI(w http.ResponseWriter, r *http.Request) {
	var igcUrl Url
	// If sent data is actual json
	if err := json.NewDecoder(r.Body).Decode(&igcUrl); err != nil {
		http.Error(w, "Check body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	track, err := igc.ParseLocation(igcUrl.Url)
	// Checks for valid URL sent in body
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {

		// Finds total track_length
		totalDistance := 0.0
		for i := 0; i < len(track.Points)-1; i++ {
			totalDistance += track.Points[i].Distance(track.Points[i+1])
		}

		// SLICE OF INT TO KEEP TRACK OF THE POST ID'S
		idCount++
		ids = append(ids, idCount)
		// Converts the id counter to string
		trackId := strconv.Itoa(idCount)
		// Stores the received track in the map
		tracks[trackId] = Tracks{
			track.Date.String(),
			track.Pilot,
			track.GliderType,
			track.GliderID,
			totalDistance,
		}
		w.Header().Set("Content-Type", "application/json")
		// Encodes unique id in json - back to user
		if err := json.NewEncoder(w).Encode(TrackId{Id: idCount}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}
}

func IdHandler(w http.ResponseWriter, r *http.Request) {

	// Regex to check the URL provided by user
	check := regexp.MustCompile("^/igcinfo/api/igc/[0-9]/([a-zA-Z_]+)$")
	test := regexp.MustCompile(("^/igcinfo/api/igc/[0-9]$"))
	path := strings.Split(r.URL.Path, "/")

	for i, v := range tracks {
		// Id in map matches Id in url
		if i == path[4] {
			// Should rather check URL before checking ID
			if check.MatchString(r.URL.Path) {
				switch path[5] {
				case "pilot":
					fmt.Fprint(w, tracks[i].Pilot)
					return
				case "glider":
					fmt.Fprint(w, tracks[i].Glider)
					return
				case "glider_id":
					fmt.Fprint(w, tracks[i].GliderId)
					return
				case "track_length":
					fmt.Fprint(w, tracks[i].Track_length)
					return
				case "H_date":
					fmt.Fprint(w, tracks[i].H_date)
					return
				default:
					// Returns 404 with empty body - when not known method is provided
					http.Error(w, "", http.StatusNotFound)
					return
				}
			} else if test.MatchString(r.URL.Path) {
				CheckHandler(w, r, v)
			} else {

				http.Error(w, "",http.StatusNotFound)
			}
			return
		}
	}
	// Returns 404 with empty body - if no ID in map matches ID provided by user
	http.Error(w, "",http.StatusNotFound)
}

func CheckHandler(w http.ResponseWriter, r *http.Request, t Tracks) {
	w.Header().Set("Content-Type", "application/json")
	// Encodes information for a specific track in json back to user
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w,"Could not encode", http.StatusInternalServerError)
	}
	return
}