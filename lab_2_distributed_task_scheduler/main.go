package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	id   int
	data string
}

func worker(id int, taskChan <-chan Task, wg *sync.WaitGroup, failChan chan<- Task) {
	defer wg.Done()

	for task := range taskChan {
		fmt.Printf("Worker %d started task %d: %s\n", id, task.id, task.data)

		// simulate random failure
		if rand.Float32() < 0.3 {
			fmt.Printf("Worker %d failed on task %d\n", id, task.id)
			failChan <- task
			return
		}

		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
		fmt.Printf("Worker %d completed task %d\n", id, task.id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	tasks := []Task{
		{id: 1, data: "Task 1"},
		{id: 2, data: "Task 2"},
		{id: 3, data: "Task 3"},
		{id: 4, data: "Task 4"},
		{id: 5, data: "Task 5"},
	}

	// initalize channels
	taskChan := make(chan Task, len(tasks))
	failChan := make(chan Task, len(tasks))

	var wg sync.WaitGroup
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, taskChan, &wg, failChan)
	}

	for _, task := range tasks {
		taskChan <- task
	}

	close(taskChan)

	go func() {
		for failedTask := range failChan {
			fmt.Printf("Reassigning failed task %d\n", failedTask.id)
			wg.Add(1)
			go worker(rand.Intn(workerCount)+1, taskChan, &wg, failChan)
		}
	}()

	wg.Wait()
	close(failChan)

	fmt.Println("All tasks completed")
}
