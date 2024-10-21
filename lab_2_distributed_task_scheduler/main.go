package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work.
type Task struct {
	id   int
	data string
}

const TASK_COUNT = 5
const RETRY_LIMIT = 1

var COMPLETED_TASKS = 0

// Worker function that processes tasks. If a worker fails, the task will be sent to failChan.
func worker(id int, taskChan <-chan Task, failChan chan<- Task, wg *sync.WaitGroup, failCounter map[int]int) {
	defer wg.Done()

	for task := range taskChan {
		// fmt.Printf("Worker %d started task %d: %s\n", id, task.id, task.data)

		// Simulate random failure (30% chance of failure)
		if rand.Float32() < 0.9 {
			fmt.Printf("Worker %d failed on task %d\n", id, task.id)
			failCounter[task.id] += 1
			if failCounter[task.id] > RETRY_LIMIT {
				fmt.Printf("Terminated task %d - exceeded retry limit of %d\n", task.id, RETRY_LIMIT)
				COMPLETED_TASKS += 1
				if COMPLETED_TASKS == TASK_COUNT {
					close(failChan)
				}
			} else {
				failChan <- task
			}

			return
		}

		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
		fmt.Printf("Worker %d completed task %d\n", id, task.id)

		COMPLETED_TASKS += 1
		if COMPLETED_TASKS == TASK_COUNT {
			close(failChan)
		}
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

	taskChan := make(chan Task, len(tasks))
	failChan := make(chan Task, len(tasks))
	failCounter := make(map[int]int)
	var wg sync.WaitGroup
	workerCount := 3

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, taskChan, failChan, &wg, failCounter)
	}

	for _, task := range tasks {
		taskChan <- task
	}

	// Handle failed tasks by redistributing them
	go func() {
		for failedTask := range failChan {
			fmt.Printf("Reassigning failed task %d\n", failedTask.id)
			wg.Add(1)
			taskChan <- failedTask
			go worker(rand.Intn(workerCount)+1, taskChan, failChan, &wg, failCounter)
		}

		close(taskChan)
	}()

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("All tasks completed.")
}
