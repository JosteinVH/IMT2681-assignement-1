package main

import (
	"encoding/json"
	"fmt"
	"github.com/marni/goigc"
	"log"
	"net/http"
	"os"
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
	var test string
	s := d.Seconds()
	if int(s)/60 == 1 {
		felles = append(felles, strconv.Itoa(int(d.Minutes())))
		felles = append(felles, "M")
		s -= 60

	}
	felles = append(felles, strconv.Itoa(int(s)))
	felles = append(felles, "S")
	test = strings.Join(felles, "")
	test = strings.Trim(test, "60")
	fmt.Println(test)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	/*if len(r.URL.Path) > 2 && len(r.URL.Path) 3 {
		uptime := time.Since(t) // returns type Duration
		conver(uptime)
	}*/

}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		fmt.Println(len(r.URL.Path))+
		if r.Body == nil {
			http.Error(w, "Require parameter", http.StatusBadRequest)
			return
		} else {

			defer r.Body.Close()

			var igcUrl Url // Empty URL
			err := json.NewDecoder(r.Body).Decode(&igcUrl)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			/*u, err := url.ParseRequestURI(igcUrl.Url)
			fmt.Println(u)
			if err != nil {
				fmt.Println("FEIL")
			}*/

			// Check if URL ends with .igc
			track, err := igc.ParseLocation(igcUrl.Url)
			// Make sure the track received is not empty
			if track.Pilot != "" {
				if err != nil {
					fmt.Errorf("Problem reading the track", err)
				}

				totalDistance := 0.0 // Finds total track_length
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

			} else {
				http.Error(w, "Invalid input: Malformed URL", http.StatusBadRequest)
			}

		}

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ids); err != nil {
			http.Error(w, "Something went wrong", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Not implemented yet", http.StatusNotImplemented)
	}
}

func idHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Split(req.URL.Path, "/")
	//id, _ := req.URL.Query()["id"]

	for i, v := range tracks {
		if i == path[3] {
			if len(path) > 4 && len(path) < 6 {
				switch path[4] {
				case "pilot":
					fmt.Fprint(w, tracks[path[3]].Pilot)
					return
				case "glider":
					fmt.Fprint(w, tracks[path[3]].Glider)
					return
				case "glider_id":
					fmt.Fprint(w, tracks[path[3]].GliderId)
					return
				case "track_length":
					fmt.Fprint(w, tracks[path[3]].Track_length)
					return
				case "H_date":
					fmt.Fprint(w, tracks[path[3]].H_date)
					return
				default:
					http.Error(w, "Not implemented yet", http.StatusNotImplemented)
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(v)
				return
			}
		}
	}
	http.Error(w, "Input: Id '"+path[3]+"' not found", http.StatusBadRequest)
}

func main() {
	port := os.Getenv("PORT")
	fmt.Println(port)

	http.HandleFunc("/api/", infoHandler)
	http.HandleFunc("/api/igc", apiHandler)
	http.HandleFunc("/api/igc/", idHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
