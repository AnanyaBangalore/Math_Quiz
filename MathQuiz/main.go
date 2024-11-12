package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func problemPuller(filename string) ([]problem, error) {
	file_obj, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err.Error())
	}
	defer file_obj.Close()

	csv_reader := csv.NewReader(file_obj)
	lines, err := csv_reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading CSV: %s", err.Error())
	}

	return parseProblem(lines), nil
}

func parseProblem(lines [][]string) []problem {
	result := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		result[i] = problem{question: lines[i][0], answer: strings.TrimSpace(lines[i][1])}
	}
	return result
}

func main() {
	filename := flag.String("f", "./quiz.csv", "file which contains the questions")
	timer := flag.Int("t", 30, "this is the timer of the quiz")
	flag.Parse()

	problems, err := problemPuller(*filename)
	if err != nil {
		fmt.Printf("Something went wrong: %s\n", err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})

	if len(problems) > 5 {
		problems = problems[:5]
	}

	correctAns := 0
	timer_obj := time.NewTimer(time.Duration(*timer) * time.Second)

problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s\n", i+1, p.question)

		select {
		case <-timer_obj.C:
			fmt.Println("\nTime's up!")
			break problemLoop
		default:
			fmt.Print("Your answer: ")
			fmt.Scanln(&answer)
			answer = strings.TrimSpace(answer)
		}

		if answer == p.answer {
			correctAns++
		}
	}

	fmt.Printf("You answered %d questions correctly.\n", correctAns)
}
