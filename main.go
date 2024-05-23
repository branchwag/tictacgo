package main

import (
	"html/template"
	"net/http"
	"strconv"
)


type Game struct {
	Board [][]string
	Player string
	Winner string
}

var game = Game{
	Board: [][]string{
		{"","",""},
		{"","",""},
		{"","",""},
	},
	Player: "X",
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, game)
}

func moveHandler(w http.ResponseWriter, r *http.Request){
	rowStr := r.URL.Query().Get("row")
	colStr := r.URL.Query().Get("col")
	row, _ := strconv.Atoi(rowStr)
	col, _ := strconv.Atoi(colStr)

	if game.Board[row][col] == "" {
		game.Board[row][col] = game.Player
		if game.Player == "X" {
			game.Player = "O"
		} else {
			game.Player="X"
		}
	}
	w.Write([]byte(game.Board[row][col]))
}


func main(){
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/move", moveHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

