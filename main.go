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
var count int

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
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS dices (id INTEGER PRIMARY KEY, name TEXT, score INTEGER, count INTEGER)")

	statement.Exec()
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("room").ParseFiles("view/index.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name = r.FormValue("name")

		sql := fmt.Sprintf("SELECT score, count, id FROM dices WHERE UPPER(name) = UPPER('%s')", name)
		rows, err := database.Query(sql)
		if err != nil {
			fmt.Println(err.Error())
			score = 0
			count = 0
			id = 0
		} else {
			for rows.Next() {
				rows.Scan(&score, &count, &id)
			}
		}

		var data = map[string]string{"name": name, "score": strconv.Itoa(score), "count": strconv.Itoa(count), "id": strconv.Itoa(id)}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeIndexPost(w http.ResponseWriter, r *http.Request) {
	database, _ := sql.Open("sqlite3", "database/dice.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS dices (id INTEGER PRIMARY KEY, name TEXT, score INTEGER, count INTEGER)")

	statement.Exec()
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("view/index.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name = r.FormValue("hide_name")
		diceId, _ := strconv.Atoi(r.FormValue("dice_id"))
		diceCount, _ := strconv.Atoi(r.FormValue("count"))
		diceScore, _ := strconv.Atoi(r.FormValue("dice_score"))
		dice := rollDice()

		sum := 0
		if dice == 1 || dice == 3 || dice == 5 {
			sum += 5
		} else {
			sum -= 3
		}

		diceScore += sum
		diceCount += 1

		if diceId != 0 {
			statement, eror := database.Prepare("UPDATE dices SET score = ?, count = ? WHERE id= ?")
			if eror != nil {
				http.Error(w, eror.Error(), http.StatusInternalServerError)
				return
			} else {
				statement.Exec(diceScore, diceCount, diceId)
			}
		} else {
			statement, eror := database.Prepare("INSERT INTO dices (name, score, count) VALUES (?, ?, ?)")

			if eror != nil {
				http.Error(w, eror.Error(), http.StatusInternalServerError)
				return
			} else {
				statement.Exec(name, diceScore, diceCount)
			}
		}

		rows, err := database.Query("SELECT score FROM dices WHERE name =" + `name`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
