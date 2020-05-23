package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func getQuestionsCSV(f io.Reader) ([][]string, error) {
	r := csv.NewReader(f)
	r.FieldsPerRecord = 2
	return r.ReadAll()
}

func readAnswer() (r string, err error) {
	_, err = fmt.Scanln(&r)
	return
}

type quizzes struct {
	total   int
	correct int
	remain  int
}

func (q quizzes) PrintResults() {
	fmt.Printf("You scored %d out of %d.\n", q.correct, q.total)
}

func (q *quizzes) Check(errOnRead error, input, expected string) {
	if errOnRead == nil && strings.TrimSpace(input) == strings.TrimSpace(expected) {
		q.correct = q.correct + 1
	}
	q.remain = q.remain - 1
}

func main() {
	filename := flag.String("cvs", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	records, err := getQuestionsCSV(f)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan bool)
	q := quizzes{remain: len(records), total: len(records)}
	// start quizz
	go func() {
		for i, record := range records {
			fmt.Printf("Problem #%d: %s = ", i+1, record[0])
			ans, err := readAnswer()
			q.Check(err, ans, record[1])
		}
		c <- true
	}()

	select {
	case <-time.After(time.Second * time.Duration(*limit)):
		fmt.Println("\nTimeout, you are failed :(")
	case <-c:
		q.PrintResults()
	}

}
