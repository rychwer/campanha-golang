package main

import (
	"campanha-golang/src/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := config.GerarRouter()

	fmt.Printf("Escutando na porta %d\n", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), router))
}