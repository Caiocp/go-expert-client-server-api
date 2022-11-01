package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Dolar struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

const (
	TWO_HUNDRED_MILLISECONDS = time.Millisecond * 200
	TEN_MILLISECONDS         = time.Millisecond * 10
)

var db *sql.DB

func main() {
	database, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		panic(err)
	}
	defer database.Close()

	db = database

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT)")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", dolarHandler)
	http.ListenAndServe(":8080", nil)
}

func dolarHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, TWO_HUNDRED_MILLISECONDS)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var dolar Dolar
	err = json.NewDecoder(res.Body).Decode(&dolar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveDolarQuote(dolar)

	json.NewEncoder(w).Encode(dolar.Usdbrl)
}

func saveDolarQuote(dolar Dolar) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, TEN_MILLISECONDS)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacao (bid) VALUES (?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	insert, err := stmt.Exec(dolar.Usdbrl.Bid)
	if err != nil {
		panic(err)
	}

	fmt.Println(insert.LastInsertId())
}
