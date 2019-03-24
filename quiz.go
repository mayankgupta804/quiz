package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var csvFlag string

func init() {
	flag.StringVar(&csvFlag, "csv", "problems.csv", "choice of .csv files loaded with problems")
}

func main() {
	flag.Parse()
	filename := csvFlag
	quiz := readCSV(filename)
	startQuiz(quiz)
}

func startQuiz(quiz []Quiz) {
	correctAnswerCount := 0
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Quiz Game ")
	for i := range quiz {
		fmt.Printf(quiz[i].Question + "\t")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimRight(answer, "\r\n")
		true := strings.Compare(answer, quiz[i].Answer)
		if true == 0 {
			fmt.Printf("Correct!!\n")
			correctAnswerCount++
		} else {
			fmt.Println("Wrong!!")
		}
	}
	fmt.Printf("Total correct answers: %d\n", correctAnswerCount)
	if correctAnswerCount == len(quiz) {
		fmt.Printf("Hurray! You answered all the questions correctly. You must be a Genius!")
	} else if correctAnswerCount == 0 {
		fmt.Printf("What a dumbass!")
	}
}

func readCSV(filename string) []Quiz {
	csvFile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var quiz []Quiz
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		quiz = append(quiz, Quiz{
			Question: line[0],
			Answer:   line[1],
		})
	}
	return quiz
}

type Quiz struct {
	Question string `json:Question`
	Answer   string `json:Answer`
}
