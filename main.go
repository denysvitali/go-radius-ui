package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/denysvitali/go-radius-ui/radius/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

func main()  {

	config, err := toml.LoadFile("./config.toml")

	if err != nil {
		log.Fatal(err)
	}

	db_user := config.GetDefault("postgres.user", "radius")
	db_password := config.GetDefault("postgres.password", "")
	db_name := config.GetDefault("postgres.database", "radius")
	sslmode := config.GetDefault("postgres.sslmode", "disable")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s",
		db_user,
		db_name,
		db_password,
		sslmode)

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/api/v1/users", GetUsers)
	router.HandleFunc("/api/v1/postauth", GetPostAuth)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello :)"))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, username, attribute, op, value FROM radcheck")

	if err != nil {
		log.Fatal(err);
	}

	var entries []models.RadCheck

	for rows.Next() {
		var entry *models.RadCheck

		entry = new(models.RadCheck)
		var id int32
		var username string
		var attribute string
		var op string
		var value string

		err = rows.Scan(&id, &username, &attribute, &op, &value)

		entry.Username = username
		entry.SetId(id)
		entry.SetAttribute(attribute)
		entry.SetOp(op)
		entry.SetValue(value)

		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, *entry)
	}

	j, err := json.Marshal(entries)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(j)
}

func GetPostAuth(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, username, pass, reply, authdate FROM radpostauth ORDER BY authdate DESC")

	if err != nil {
		log.Fatal(err)
	}

	var entries []models.RadPostAuth

	for rows.Next() {
		var entry *models.RadPostAuth

		entry = new(models.RadPostAuth)
		var id int32
		var username string
		var pass string
		var authdate time.Time
		var reply string

		err = rows.Scan(&id, &username, &pass, &reply, &authdate)

		entry.Username = username
		entry.Reply = reply
		entry.AuthDate = authdate
		entry.SetId(id)
		entry.SetPass(pass)

		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, *entry)
	}

	j, err := json.Marshal(entries)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(j)
}