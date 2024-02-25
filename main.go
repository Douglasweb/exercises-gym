package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Exercise struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Created_At  time.Time `json:"datetime"`
}

var Exercises []Exercise

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllExercises(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllExercises")
	json.NewEncoder(w).Encode(Exercises)
}

func returnSingleExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	ip := ReadUserIP(r)

	fmt.Println("Ip:" + ip)

	newkey, errorParse := uuid.Parse(key)

	if errorParse != nil {
		fmt.Println("Error: " + errorParse.Error())
		json.NewEncoder(w).Encode(errorParse.Error())
		return
	}

	response, err := GetSQLByID(newkey)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		json.NewEncoder(w).Encode(nil)
		return
	} else {
		json.NewEncoder(w).Encode(response)
		return
	}

}

func createNewExercise(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := io.ReadAll(r.Body)
	var exercise Exercise
	json.Unmarshal(reqBody, &exercise)

	json.NewEncoder(w).Encode(exercise)
}

func deleteExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/exercises", returnAllExercises)
	myRouter.HandleFunc("/exercise", createNewExercise).Methods("POST")
	myRouter.HandleFunc("/exercise/{id}", deleteExercise).Methods("DELETE")
	myRouter.HandleFunc("/exercise/{id}", returnSingleExercise)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func OpenConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=teste password=teste dbname=teste sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	return db, err
}

func Insert(e Exercise) (id string, err error) {
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := `INSERT INTO  maniac.exercises (nome) VALUES ($1) RETURNING id`
	err = conn.QueryRow(sql, e.Name).Scan(&id)
	return id, err
}

func GetSQLByID(idKey uuid.UUID) (e Exercise, err error) {
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow(`SELECT id, name, description, created_at FROM maniac.exercises WHERE id=$1`, idKey)

	err = row.Scan(&e.Id, &e.Name, &e.Description, &e.Created_At)
	return
}

func main() {

	handleRequests()
}
