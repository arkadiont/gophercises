package main

import (
	"flag"
	"fmt"
	cyoa "gophercises/ex3-cyoa"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8080, "specify listening port")
	file := flag.String("file", "gopher.json", "a JSON file with CYOA story")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	b, err := cyoa.GetBook(f)
	if err != nil {
		log.Fatal(err)
	}

	listen := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting the server on %s\n", listen)
	log.Fatal(http.ListenAndServe(listen, cyoa.NewHandler(b)))
}
