package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Dolar struct {
	Bid string `json:"bid"`
}

const (
	THREE_HUNDRED_MILLISECONDS = time.Millisecond * 300
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, THREE_HUNDRED_MILLISECONDS)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
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
		panic(err)
	}

	file, err := os.Create("./cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", dolar.Bid))
}
