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

// Total count of tasks
const TASK_COUNT = 5

// Limit of the number of times a process can be retried
const RETRY_LIMIT = 1

// Tasks currently completed
var COMPLETED_TASKS = 0

// Worker function that processes tasks. If a worker fails, the task will be sent to failChan.
func worker(id int, taskChan <-chan Task, failChan chan<- Task, wg *sync.WaitGroup, failCounter map[int]int) {
	defer wg.Done()

	for task := range taskChan {
		fmt.Printf("\t - Worker %d started task %d: %s\n", id, task.id, task.data)

		// Simulate random failure (30% chance of failure)
		if rand.Float32() < 0.3 {
			fmt.Printf("\033[31m\t - Worker %d failed on task %d\033[0m\n", id, task.id)
			failCounter[task.id] += 1

			if failCounter[task.id] > RETRY_LIMIT {
				fmt.Printf("\033[33m\t - Terminated task %d - exceeded retry limit of %d\033[0m\n", task.id, RETRY_LIMIT)

				// Close the failed task channel if all the tasks have been completed
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
		fmt.Printf("\033[32m\t - Worker %d completed task %d\033[0m\n", id, task.id)

		// Close the failed task channel if all the tasks have been completed
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

	// channels
	taskChan := make(chan Task, len(tasks))
	failChan := make(chan Task, len(tasks))

	// map for keeping track of the failures
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
		currentWorkerID := 0
		for failedTask := range failChan {
			fmt.Printf("\033[37m\t - Reassigning failed task %d to worker %d\033[0m\n", failedTask.id, currentWorkerID+1)
			wg.Add(1)
			taskChan <- failedTask
			go worker(currentWorkerID+1, taskChan, failChan, &wg, failCounter)

			currentWorkerID = (currentWorkerID + 1) % workerCount // cycle worker id
		}

		close(taskChan)
	}()

	// Wait for all workers to finish
	wg.Wait()
	time.Sleep(time.Second)

	fmt.Println("\n\tAll tasks completed")
}
