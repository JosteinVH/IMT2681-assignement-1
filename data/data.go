package data

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