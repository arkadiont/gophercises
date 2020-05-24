package main

import (
	"fmt"
	link "gophercises/ex4-link"
	"log"
	"os"
)

var htmlEx1 = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
  <a href="/page-two">A link 
  to a second
   page</a>
</body>
</html>
`

func main() {
	f, err := os.Open("../ex4.html")
	if err != nil {
		log.Fatal(err)
	}
	links, err := link.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("links: %d\n%+v\n", len(links), links)
}
