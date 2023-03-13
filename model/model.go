package model

type Movie struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Genre  string `json:"genre" db:"genre"`
	Rating int    `json:"rating" db:"rating"`
}

type MovieFile struct {
	Name   string `json:"name" db:"name"`
	Genre  string `json:"genre" db:"genre"`
	Rating int    `json:"rating" db:"rating"`
}
