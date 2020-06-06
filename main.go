package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var name string
var id int
var score int

var onlyOnce sync.Once

// prepare the dice
var dice = []int{1, 2, 3, 4, 5, 6}

func rollDice() int {

	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano()) // only run once
	})

	return dice[rand.Intn(len(dice))]
}

func main() {

	http.HandleFunc("/", routeIndexGet)
	http.HandleFunc("/room", routeRoomPost)
	http.HandleFunc("/process", routeIndexPost)

	fmt.Println("server started at localhost:9001")
	http.ListenAndServe(":9001", nil)
	// name := name()

	// fmt.Println("Name ", name)

	// dice1 := rollDice()

	// fmt.Println("Dice 1: ", dice1)
}

func routeIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.New("form").ParseFiles("view/index.html"))
		err := tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeRoomPost(w http.ResponseWriter, r *http.Request) {
	database, _ := sql.Open("sqlite3", "database/dice.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS dices (id INTEGER PRIMARY KEY, name TEXT, score INTEGER)")

	statement.Exec()
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("room").ParseFiles("view/index.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name = r.FormValue("name")

		rows, err := database.Query("SELECT SUM(score) as score FROM dices WHERE name =" + name)
		if err != nil {
			score = 0
		} else {
			for rows.Next() {
				rows.Scan(&score)
			}
		}

		var data = map[string]string{"name": name, "score": strconv.Itoa(score)}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeIndexPost(w http.ResponseWriter, r *http.Request) {
	database, _ := sql.Open("sqlite3", "database/dice.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS dices (id INTEGER PRIMARY KEY, name TEXT, score INTEGER)")

	statement.Exec()
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("view/index.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name = r.FormValue("hide_name")
		dice := rollDice()

		sum := 0
		if sum == 1 || sum == 3 || sum == 5 {
			sum += 5
		} else {
			sum -= 3
		}

		statement, _ = database.Prepare("INSERT INTO dices (name, score) VALUES (?, ?)")
		statement.Exec(name, sum)

		rows, err := database.Query("SELECT SUM(score) as score FROM dices WHERE name =" + name)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
		} else {
			for rows.Next() {
				rows.Scan(&score)
			}
		}

		var data = map[string]string{"name": name, "dice": strconv.Itoa(dice), "score": strconv.Itoa(score)}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
