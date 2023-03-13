package dbHelper

import (
	"z_test/db"
	"z_test/model"
)

var database = db.DBConnect()

func CreateMovie(item model.Movie) error {
	// language = SQL
	query := "insert into movies(name, genre, rating) values ($1,$2,$3)"
	_, err := database.Exec(query, item.Name, item.Genre, item.Rating)
	return err
}

func ReadMovie() ([]model.Movie, error) {
	// language = SQL
	query := "select * from movies"
	list := make([]model.Movie, 0)
	err := database.Select(&list, query)
	return list, err
}

func UpdateMovie(item model.Movie) error {
	// language = SQL
	query := "update movies set name=$1, genre=$2, rating=$3 where id=$4"
	_, err := database.Exec(query, item.Name, item.Genre, item.Rating, item.ID)
	return err
}

func DeleteMovie(id int) error {
	// language = SQL
	query := "delete from movies where id=$1"
	_, err := database.Exec(query, id)
	return err
}
