package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	BASE_URL = "http://localhost:8080/cotacao"
)

type QuoteHTTPResponse struct {
	Dib string `json:"bid"`
}

func main() {
	contextQuote := context.Background()
	contextQuote, cancel := context.WithTimeout(contextQuote, 1000*time.Millisecond)
	defer cancel()
	request, errorRequest := http.NewRequestWithContext(contextQuote, http.MethodGet, BASE_URL, nil)
	if errorRequest != nil {
		panic(errorRequest)
	}
	output, errorOutput := http.DefaultClient.Do(request)
	if errorOutput != nil {
		log.Fatal(errorOutput)
	}
	defer output.Body.Close()
	body, errorBody := io.ReadAll(output.Body)
	if errorBody != nil {
		panic(errorBody)
	}
	var quoteResponse QuoteHTTPResponse
	errorUnmarshal := json.Unmarshal(body, &quoteResponse)
	if errorUnmarshal != nil {
		panic(errorUnmarshal)
	}
	var file *os.File
	defer file.Close()
	_, errorCreateFile := os.Stat("quote.txt")
	if os.IsNotExist(errorCreateFile) {
		file, errorCreateFile = os.Create("quote.txt")
		if errorCreateFile != nil {
			panic(errorCreateFile)
		}
	} else {
		file, errorCreateFile = os.OpenFile("quote.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if errorCreateFile != nil {
			log.Fatal(errorCreateFile)
			panic(errorCreateFile)
		}
	}
	currentTime := time.Now()
	currentTimeFormatted := currentTime.Format("2006-01-02 15:04:05")
	_, errorWriteFile := file.WriteString(currentTimeFormatted + " - " + quoteResponse.Dib + "\n")
	if errorWriteFile != nil {
		panic(errorWriteFile)
	}
}
