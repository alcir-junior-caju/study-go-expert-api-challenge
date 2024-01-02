package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	SERVER_PORT = "8080"
	BASE_URL    = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	MYSQL       = "yzwudqvd0vw7h5dwf3bs:pscale_pw_gBCsRDvPaZ9zFsgzSV1ej04BMO47LtpL283nwv8whF7@tcp(aws.connect.psdb.cloud)/go?tls=true&interpolateParams=true"
)

type Quote struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type QuoteHTTPResponse struct {
	USDBRL Quote `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", quoteHandler)
	fmt.Println("Server API listening on port", SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+SERVER_PORT, nil))
}

func getQuote() (*Quote, error) {
	contextQuote := context.Background()
	contextQuote, cancel := context.WithTimeout(contextQuote, 300*time.Millisecond)
	defer cancel()
	responseQuote := QuoteHTTPResponse{}
	requestQuote, errorRequestQuote := http.NewRequestWithContext(contextQuote, http.MethodGet, BASE_URL, nil)
	if errorRequestQuote != nil {
		panic(errorRequestQuote)
	}
	outputQuote, errorOutputQuote := http.DefaultClient.Do(requestQuote)
	if errors.Is(errorOutputQuote, context.DeadlineExceeded) {
		log.Println("Context Error:", errorOutputQuote)
	}
	if errorOutputQuote != nil {
		panic(errorOutputQuote)
	}
	defer outputQuote.Body.Close()
	errorOutputQuote = json.NewDecoder(outputQuote.Body).Decode(&responseQuote)
	if errorOutputQuote != nil {
		panic(errorOutputQuote)
	}
	return &responseQuote.USDBRL, nil
}

func quoteHandler(writer http.ResponseWriter, _ *http.Request) {
	quoteResponse, err := getQuote()
	if err != nil {
		http.Error(writer, "Error fetching quote", http.StatusInternalServerError)
		return
	}
	database, errorDatabase := connectionDatabase()
	if errorDatabase != nil {
		http.Error(writer, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer database.Close()
	errorInsertQuote := insertQuote(database, quoteResponse)
	if errorInsertQuote != nil {
		http.Error(writer, "Error inserting quote", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(quoteResponse)
}

func connectionDatabase() (*sql.DB, error) {
	database, errorDatabase := sql.Open("mysql", MYSQL)
	if errorDatabase != nil {
		log.Fatalf("failed to connect: %v", errorDatabase)
	}
	if errorDatabasePing := database.Ping(); errorDatabasePing != nil {
		log.Fatalf("failed to ping: %v", errorDatabasePing)
	}
	log.Println("Successfully connected to PlanetScale!")
	// TODO: Create initial table PlanetScale
	statement, errorStatement := database.Prepare("CREATE TABLE IF NOT EXISTS quotes (id INT AUTO_INCREMENT PRIMARY KEY, code VARCHAR(255), codein VARCHAR(255), name VARCHAR(255), high VARCHAR(255), low VARCHAR(255), varBid VARCHAR(255), pctChange VARCHAR(255), bid VARCHAR(255), ask VARCHAR(255), timestamp VARCHAR(255), createDate VARCHAR(255))")
	if errorStatement != nil {
		return nil, errorStatement
	}
	defer statement.Close()
	_, errorExec := statement.Exec()
	if errorExec != nil {
		return nil, errorExec
	}
	return database, nil
}

func insertQuote(database *sql.DB, quote *Quote) error {
	contextDatabaseInsert := context.Background()
	contextDatabaseInsert, cancel := context.WithTimeout(contextDatabaseInsert, 400*time.Millisecond)
	defer cancel()
	statement, errorStatement := database.Prepare("INSERT INTO quotes (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if errorStatement != nil {
		log.Printf("Error preparing SQL: %v", errorStatement)
		return errorStatement
	}
	defer statement.Close()
	_, errorExec := statement.ExecContext(contextDatabaseInsert, quote.Code, quote.Codein, quote.Name, quote.High, quote.Low, quote.VarBid, quote.PctChange, quote.Bid, quote.Ask, quote.Timestamp, quote.CreateDate)
	if errorExec != nil {
		log.Printf("Error executing SQL: %v", errorExec)
		log.Printf("Generated SQL: %v", statement)
		return errorExec
	}
	defer database.Close()
	return nil
}
