package api

import (
	"encoding/json"
	"fmt"
	."jvh_local/TEST/data"
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
	/*min := int(sec) / 60
	hour := 1 //int(sec) / 3600
	day := int(sec) / 86400
	month := int(sec) / 2629746
	year := int(sec) / 31556952*/
const (
	multMinute = 60
	multHour   = multMinute * 60
	multDay    = multHour * 24
	multWeek   = multDay * 7
	multYear   = multDay * 365.25 // have to get years first to account for leap years
	multMonth  = multYear / 12    // once we have years we use that to get months
)
	felles = append(felles, "P")

	year := int(sec / multYear)
	if year >= 1 {
		felles = append(felles, strconv.Itoa(year))
		felles = append(felles, "Y")
		sec -= float64(31556952 * year)
	}

	month := int(sec / multMonth)
	if month >= 1 {
		felles = append(felles, strconv.Itoa(month))
		felles = append(felles, "M")
		sec -= float64(2629746 * month)
	}

	day := int(sec / multDay)
	if day >= 1 {
		felles = append(felles, strconv.Itoa(day))
		felles = append(felles, "D")
		sec -= float64(86400 * day)
	}

	felles = append(felles, "T")

	hour := int(sec / multHour)
	if hour >= 1 {
		felles = append(felles, strconv.Itoa(hour))
		felles = append(felles, "H")
		sec -= float64(3600 * hour)

	}

	min := int(sec / multMinute)
	if min >= 1 {
		felles = append(felles, strconv.Itoa(min))
		felles = append(felles, "M")
		sec -= float64(60 * min)

	}
	if sec >= 0 {
		felles = append(felles, strconv.Itoa(int(sec)))
		felles = append(felles, "S")
	}

	k := strings.Join(felles, "")
	return k
}



func InfoHandler(w http.ResponseWriter, r *http.Request) {
	check := regexp.MustCompile("^/igcinfo/api/$")
	if check.MatchString(r.URL.Path) {
		// Time since application started
		uptime := time.Since(t)
		iso := conver(uptime)
		infoApi := Info{
			iso,
			"Service for IGC tracks.",
			"v1",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(infoApi)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
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
	if err := json.NewDecoder(r.Body).Decode(&igcUrl); err != nil {
		http.Error(w, "Check body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	track, err := igc.ParseLocation(igcUrl.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func GetAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}
}

func IdHandler(w http.ResponseWriter, r *http.Request) {

	check := regexp.MustCompile("^/igcinfo/api/igc/[0-9]/([a-zA-Z_]+)$")
	test := regexp.MustCompile(("^/igcinfo/api/igc/[0-9]$"))
	path := strings.Split(r.URL.Path, "/")

	for i, v := range tracks {
		if i == path[4] {
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
					http.Error(w, "", http.StatusNotFound)
					return
				}
			} else if test.MatchString(r.URL.Path) {
				checkHandler(w, r, v)
			} else {
				http.NotFound(w, r)
			}
			return
		}
	}
	http.Error(w, "",http.StatusNotFound)
	//http.NotFound(w,r)
}

func checkHandler(w http.ResponseWriter, r *http.Request, t Tracks) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w,"Could not encode", http.StatusInternalServerError)
	}
	return
}