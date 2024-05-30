package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Game struct {
	Board  [][]string
	Player string
}

var game = Game{
	Board: [][]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	},
	Player: "X",
}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("index.html", "board.html"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index page")
	err := templates.ExecuteTemplate(w, "index.html", game)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received move request")
	rowStr := r.URL.Query().Get("row")
	colStr := r.URL.Query().Get("col")
	log.Printf("Move request parameters: row=%s, col=%s", rowStr, colStr)

	row, errRow := strconv.Atoi(rowStr)
	col, errCol := strconv.Atoi(colStr)

	if errRow != nil || errCol != nil || row < 0 || row >= len(game.Board) || col < 0 || col >= len(game.Board[0]) {
		http.Error(w, "Invalid row or column", http.StatusBadRequest)
		return
	}

	if game.Board[row][col] == "" {
		game.Board[row][col] = "X"
		togglePlayer()
		computerMove() 
	} else {
		http.Error(w, "Invalid move, cell already occupied", http.StatusBadRequest)
		return
	}

	err := templates.ExecuteTemplate(w, "board.html", game)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
	log.Println("Move processed and board updated")
}

func togglePlayer() {
	if game.Player == "X" {
		game.Player = "O"
	} else {
		game.Player = "X"
	}
}

func computerMove() {
	for i, row := range game.Board{
		for j, cell := range row {
			if cell == "" {
				game.Board[i][j] = "O"
				togglePlayer()
				return
			}
		}
	}
}

func main() {
	port := ":8080"
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/move", moveHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	log.Printf("Serving it up on http://localhost%s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
