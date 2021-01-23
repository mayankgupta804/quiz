package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var filename string
var timeLimit int
var shuffle bool

func init() {
	flag.StringVar(&filename, "csv", "problems.csv", "choice of .csv files loaded with problems")
	flag.IntVar(&timeLimit, "time", 30, "time limit for the quiz (in number of seconds)")
	flag.BoolVar(&shuffle, "shuffle", false, "to shuffle or not? that is the question. (options=true/false)")
}

func main() {
	flag.Parse()
	quiz := readCSV(filename)
	if shuffle {
		shuffleQuiz(quiz)
	}
	initializeQuiz()
	fmt.Printf("Time limit for the quiz is: %d\n", timeLimit)
	done := make(chan struct{})
	time.AfterFunc(time.Duration(timeLimit)*time.Duration(time.Second), func() {
		done <- struct{}{}
	})
	startQuiz(quiz, done)
}

func shuffleQuiz(quiz []Quiz) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quiz), func(i, j int) { quiz[i], quiz[j] = quiz[j], quiz[i] })
}

func initializeQuiz() {
	fmt.Println("Welcome to the Quiz Game")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press ENTER to start")
	reader.ReadByte()
}

func startQuiz(quiz []Quiz, done <-chan struct{}) {
	correctAnswerCount := 0
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for i := range quiz {
			fmt.Printf(quiz[i].Question + "\t")
			input, _ := reader.ReadString('\n')
			trimmedInput := strings.TrimRight(input, "\r\n")
			isCorrect := strings.Compare(trimmedInput, quiz[i].Answer)
			if isCorrect == 0 {
				correctAnswerCount++
			}
		}
	}()
	<-done
	fmt.Printf("\nTotal correct answers: %d\n", correctAnswerCount)
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
