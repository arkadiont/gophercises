package main

import (
	"flag"
	"fmt"
	urlshort "gophercises/ex2-urlshort"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "specify listening port")

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	// Json
	json := `[{
		"path": "/urlshortV2",
		"url": "https://github.com/gophercises/urlshort"
	}]`

	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	listen := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting the server on %s\n", listen)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", notFound)
	return mux
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "Not found!")
}
