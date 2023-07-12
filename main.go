package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// currently no database is used, so the data will be inserted in the slice
var movies []Movie

// struct for content of a movie
type Movie struct {
	// `` : is used to show the specified string to the client
	//client will see "id" not the "ID"
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`

	// values in the Director struct will be associated here
	Director *Director `json:"director"`
}

// movie and director are associated
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// functions
func getMovies(w http.ResponseWriter, r *http.Request) {
	// set the content type of our response to the json
	w.Header().Set("Content-Type", "application/json")

	// create a new json encoder for the response w, and encode the current slice movies
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//get the movie id from the request to delete the movie
	params := mux.Vars(r)

	//search through the slice for the movie id
	for index, item := range movies {
		if item.ID == params["id"] {
			//delete the match movie
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	//give the remaining movies
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//search for the movie
	for _, item := range movies {
		if item.ID == params["id"] {
			//encode the found item for the response
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//create a temp movie to append
	var movie Movie

	// decode the incoming request body into the temp movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	//generate a random movie ID for the new movie
	movie.ID = strconv.Itoa(rand.Intn(100000000))

	//append the movie
	movies = append(movies, movie)

	//return the create movie as success
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			//delete the matched movie
			movies = append(movies[:index], movies[index+1:]...)

			//create a new temp movie to append
			var movie Movie
			//decode the current request body into temp movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			//set the movie id to previous movie ID
			movie.ID = params["id"]
			//append the new movie to current slice
			movies = append(movies, movie)
			//return the added movie
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {

	//append some movies to the slice
	movies = append(movies, Movie{ID: "1", Isbn: "100101", Title: "Movie one", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "101101", Title: "Movie two", Director: &Director{Firstname: "Akhil", Lastname: "Sharma"}})

	// create a router to handle the routes
	r := mux.NewRouter()

	// create routes to handle the functions
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET") //{id} will recieve the id
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("starting server at port :8000 ...\n")

	//start the server
	log.Fatal(http.ListenAndServe(":8000", r))
}
