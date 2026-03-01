package seeder

import "time"

type Game struct {
	Date    time.Time `json:"date"`
	Format  string    `json:"format"`
	Results []Result  `json:"results"`
}

type Result struct {
	Player string `json:"player"`
	Name   string `json:"name"`
	Place  int    `json:"place"`
	Kills  int    `json:"kills"`
}
