package seeder

import "time"

type Game struct {
	Date    time.Time `json:"date"`
	Results []Result  `json:"results"`
}

type Result struct {
	Player    string `json:"player"`
	Commander string `json:"commander"`
	Place     int    `json:"place"`
	Kills     int    `json:"kills"`
}
