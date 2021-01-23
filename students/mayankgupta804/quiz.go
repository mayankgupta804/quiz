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
	"sync"
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
	quiz, err := readCSVFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if shuffle {
		shuffleQuiz(quiz)
	}
	initializeQuiz()
	fmt.Printf("\t\t\t\tTime limit for the quiz is: %d seconds\n", timeLimit)
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
	fmt.Println("\t\t\t\tWelcome to the Quiz")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\t\t\t\tPress ENTER to start")
	reader.ReadByte()
}

func startQuiz(quiz []Quiz, done <-chan struct{}) {
	correctAnswerCount := 0
	reader := bufio.NewReader(os.Stdin)
	mutex := &sync.Mutex{}
	go func(mutex *sync.Mutex) {
		for i := range quiz {
			fmt.Printf(quiz[i].Question + "\t")
			input, _ := reader.ReadString('\n')
			trimmedInput := strings.TrimRight(input, "\r\n")
			isCorrect := strings.Compare(trimmedInput, quiz[i].Answer)
			if isCorrect == 0 {
				mutex.Lock()
				correctAnswerCount++
				mutex.Unlock()
			}
		}
	}(mutex)
	<-done
	mutex.Lock()
	fmt.Printf("\nTotal correct answers: %d/%d\n", correctAnswerCount, len(quiz))
	mutex.Unlock()
	if correctAnswerCount == len(quiz) {
		fmt.Printf("Hurray! You answered all the questions correctly.")
	} else if correctAnswerCount == 0 {
		fmt.Printf("Ouch!")
	}
}

func readCSVFile(filename string) ([]Quiz, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var quiz []Quiz
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		quiz = append(quiz, Quiz{
			Question: line[0],
			Answer:   line[1],
		})
	}
	return quiz, nil
}

// Quiz holds a single question and its respective answer
type Quiz struct {
	Question string
	Answer   string
}
