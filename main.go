package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":7171"

	http.HandleFunc("/search", getScrapedData)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
