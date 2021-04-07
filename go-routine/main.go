package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Todo struct {
	UserID    int64  `json:"userId"`
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func getFromEndpoint(endpoint string, c chan *Todo) (*Todo, error) {
	var todo Todo
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Println(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&todo)
	if c == nil {
		if err != nil {
			log.Println(err)
		}
		if err != nil {
			return nil, err
		}
	} else {
		c <- &todo
	}
	return &todo, nil
}

func GetSequential(endpoints []string) []Todo {
	var todos []Todo
	for _, value := range endpoints {

		todo, _ := getFromEndpoint(value, nil)
		todos = append(todos, *todo)
	}
	return todos
}

func GetGoroutine(endpoints []string) (todos []Todo) {
	var wg sync.WaitGroup
	before := time.Now()
	c := make(chan *Todo, 5)
	for _, value := range endpoints {
		wg.Add(1)
		go func(x string) {
			defer wg.Done()
			_, _ = getFromEndpoint(x, c)
		}(value)
	}

	go func() {
		wg.Wait()
		close(c)
		defer log.Println(time.Since(before).Milliseconds())
	}()

	for value := range c {
		todos = append(todos, *value)
	}

	return todos
}

func main() {

	endpoints := []string{
		"https://jsonplaceholder.typicode.com/todos/1",
		"https://jsonplaceholder.typicode.com/todos/2",
		"https://jsonplaceholder.typicode.com/todos/3",
		"https://jsonplaceholder.typicode.com/todos/4",
		"https://jsonplaceholder.typicode.com/todos/5",
	}

	// without goroutine
	timeBefore := time.Now()
	todos := GetSequential(endpoints)
	log.Println(todos)
	log.Println(time.Since(timeBefore).Milliseconds())

	// with goroutine
	todos2 := GetGoroutine(endpoints)
	log.Println(todos2)
}
