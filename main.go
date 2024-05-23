package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var tmpl = template.Must(template.ParseFiles("index.html"))
var mu sync.Mutex

type Game struct {
	Board [3][3]string
	Player string
	Winner string
}

var game = Game{
	Board: [3][3]string{{"","",""}, {"","",""}, {"","",""}},
	Player: "X",
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	mu.Lock()
	defer mu.Unlock()
	tmpl.Execute(w, game)
}

func moveHandler(w http.ResponseWriter, r *http.Request){
	mu.Lock()
	defer mu.Unlock()

	if game.Winner != "" {
		tmpl.Execute(w, game)
		return
	}

	row, err1 := strconv.Atoi(r.URL.Query().Get("row"))
	col, err2 := strconv.Atoi(r.URL.Query().Get("col"))

	if err1 != nil || err2 != nil || row < 0 || row > 2 || col < 0 || col > 2{
		tmpl.Execute(w, game)
		return
	}

	if game.Board[row][col] == "" {
		game.Board[row][col] = game.Player
		if checkWinner(game.Board, game.Player){
			game.Winner = game.Player
		} else {
			if game.Player == "X" {
				game.Player = "O"
			} else {
				game.Player = "X"
			}
		}
	}

	tmpl.Execute(w, game)
}

func checkWinner(board [3][3]string, player string) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] ==player {
		return true
	}
	return false
}

func main(){
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/move", moveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

