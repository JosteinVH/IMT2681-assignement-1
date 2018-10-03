package main

import (
	"encoding/json"
	"fmt"
	"github.com/marni/goigc"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Tracks struct {
	H_date       string  `json:"H_date"`
	Pilot        string  `json:"pilot"`
	Glider       string  `json:"glider"`
	GliderId     string  `json:"glider_id"`
	Track_length float64 `json:"track_length"`
}

type Info struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Track ids
type TrackId struct {
	Id int `json:"id"`
}

// POST URL
type Url struct {
	Url string `json:"url"`
}

var idCount int = 0
var ids = make([]int, 0)
var tracks = make(map[string]Tracks)
var t = time.Now()

func conver(d time.Duration) {
	var felles = make([]string, 0)
	sec := d.Seconds()
	min := int(sec) / 60
	hour := int(sec) / 3600
	day := int(sec) / 86400
	month := int(sec) / 2629746
	year := int(sec) / 31556952

	felles = append(felles, "P")
	if year >= 1 {
		felles = append(felles, "Y")
		felles = append(felles, strconv.Itoa(year))

		sec -= float64(31556952 * year)
	}

	if month >= 1 {
		felles = append(felles, "M")
		felles = append(felles, strconv.Itoa(month))

		sec -= float64(2629746 * month)
	}

	if day >= 1 {
		felles = append(felles, "D")
		felles = append(felles, strconv.Itoa(day))

		sec -= float64(86400 * day)
	}

	if hour >= 1 {
		felles = append(felles, "H")
		felles = append(felles, strconv.Itoa(hour))

		sec -= float64(3600 * hour)
	}
	if min >= 1 {
		felles = append(felles, "M")
		felles = append(felles, strconv.Itoa(min))

		sec -= float64(60 * sec)
	}
	if sec > 0 {
		felles = append(felles, strconv.Itoa(int(sec)))
		felles = append(felles, "S")
	}

	fmt.Println(felles)
	k := strings.Join(felles, "")
	fmt.Println(k)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kjøres")
	check := regexp.MustCompile("^/api/$")
	if check.FindString(r.URL.Path) == "/api/" {
		uptime := time.Since(t) // returns type Duration
		conver(uptime)

	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	//check := regexp.MustCompile("^/api/igc")

	switch r.Method {
	case "POST":
		//function
		postAPI(w, r)

	case "GET":
		getAPI(w, r)
	default:
		http.Error(w, "Not implemented yet", http.StatusNotImplemented)
	}

}

func postAPI(w http.ResponseWriter, r *http.Request) {
	/*if r.Body == nil {
		http.Error(w, "Require parameter", http.StatusBadRequest)
		return
	} else {*/
	// Regexp for matching url.
	/*check, _ := regexp.Compile("^http://skypolaris.org/wp-content/uploads/IGS([a-zA-Z0-9/.%-()=!#¤%&\"]+)Files/([a-zA-Z0-9/.%-()=!#¤%&]+).igc$")
	check.MatchString(igcUrl.Url)*/
	// Empty URL
	var igcUrl Url
	if err := json.NewDecoder(r.Body).Decode(&igcUrl); err != nil {
		http.Error(w, "Check body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	track, err := igc.ParseLocation(igcUrl.Url)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
		http.Error(w, "No data", http.StatusBadRequest)
		return
	}

	// Make sure the track received is not empty
	// Finds total track_length
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	// SLICE OF INT TO KEEP TRACK OF THE POST ID'S
	idCount++
	ids = append(ids, idCount)
	// Converts the id coutner to string
	trackId := strconv.Itoa(idCount)
	tracks[trackId] = Tracks{
		track.Date.String(),
		track.Pilot,
		track.GliderType,
		track.GliderID,
		totalDistance,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(TrackId{Id: idCount}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}
}

func idHandler(w http.ResponseWriter, r *http.Request) {

	check := regexp.MustCompile("^/api/igc/[0-9]/([a-zA-Z_]+)$")
	test := regexp.MustCompile(("^/api/igc/[0-9]$"))
	path := strings.Split(r.URL.Path, "/")

	for i, v := range tracks {
		if i == path[3] {
			if check.MatchString(r.URL.Path) {
				switch path[4] {
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
					/*default:
					http.Error(w, "Missing", http.StatusNotFound)
					return*/
				}
			} else if test.MatchString(r.URL.Path) {
				checkHandler(w, r, v)
			}
		}
	}
	http.Error(w, "Input: Id '"+path[3]+"' not found", http.StatusNotFound)
}
func metaHandler(w http.ResponseWriter, r *http.Request, k string) {

}

func checkHandler(w http.ResponseWriter, r *http.Request, k Tracks) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(k)
	return
}

func main() {
	http.HandleFunc("/api/", infoHandler)
	http.HandleFunc("/api/igc", apiHandler)
	http.HandleFunc("/api/igc/", idHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
