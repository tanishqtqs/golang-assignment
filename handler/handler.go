package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"z_test/db/dbHelper"
	"z_test/model"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body := "Hello World"
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func AddMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	createErr := dbHelper.CreateMovie(movie)
	if createErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(createErr)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("Movie added successfully")
	if err != nil {
		logrus.Error(err)
		return
	}
}

func GetMovie(w http.ResponseWriter, r *http.Request) {

	movies, err := dbHelper.ReadMovie()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	updErr := dbHelper.UpdateMovie(movie)
	if updErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(updErr)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("Movie updated successfully")
	if err != nil {
		logrus.Error(err)
		return
	}
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	delErr := dbHelper.DeleteMovie(id)
	if delErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(delErr)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("Movie deleted successfully")
	if err != nil {
		logrus.Error(err)
		return
	}
}

func ReadCSV(w http.ResponseWriter, r *http.Request) {
	// input file from multipart form body
	file, _, _ := r.FormFile("csvFile")

	// create temporary csv file
	create, createErr := os.Create("temp.csv")
	if createErr != nil {
		log.Print(createErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// copy file contents into our temp.csv
	_, copyErr := io.Copy(create, file)
	if copyErr != nil {
		log.Print(copyErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Open("./temp.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	data, readErr := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(readErr)
	}
	// in our example we assume we are getting a list of movies from the file
	movieList := make([]model.MovieFile, 0)
	for i, line := range data {
		if i == 0 {
			continue
		}
		var movie model.MovieFile
		movie.Name = line[0]
		movie.Genre = line[1]
		rating := line[2]
		movie.Rating, err = strconv.Atoi(rating)

		movieList = append(movieList, movie)
	}

	remErr := os.Remove("./temp.csv")
	if remErr != nil {
		log.Print(remErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(movieList)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	filepath := path.Base(resp.Request.URL.String())
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return filepath, err
}

func DownloadMultipleFiles(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("urls")
	// In our example we get input urls in url params in the form : abc.com,123.com,xyz.com
	urls := strings.Split(urlStr, ",")
	done := make(chan string, len(urls))
	errChan := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := DownloadFile(URL)
			if err != nil {
				errChan <- err
				done <- ""
				return
			}
			done <- b
			errChan <- nil
		}(URL)
	}
	downloadedFilePaths := make([]string, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {
		downloadedFilePaths = append(downloadedFilePaths, <-done)
		if err := <-errChan; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(downloadedFilePaths)
	if err != nil {
		logrus.Error(err)
		return
	}
}
