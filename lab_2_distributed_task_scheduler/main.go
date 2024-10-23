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

const TASK_COUNT = 5
const RETRY_LIMIT = 1

var COMPLETED_TASKS = 0

func worker(id int, taskChan <-chan Task, failChan chan<- Task, wg *sync.WaitGroup, failCounter map[int]int) {
	defer wg.Done()

	for task := range taskChan {
		fmt.Printf(" | Worker %d started task %d: %s\n", id, task.id, task.data)

		// simulate random failure
		if rand.Float32() < 0.35 {
			fmt.Printf("\033[31m | Worker %d failed on task %d\033[0m\n", id, task.id)
			failCounter[task.id] += 1

			if failCounter[task.id] > RETRY_LIMIT {
				fmt.Printf("\033[33m | Terminated task %d - exceeded retry limit of %d\033[0m\n", task.id, RETRY_LIMIT)

				// close the failed task channel if all the tasks have been completed
				COMPLETED_TASKS += 1
				if COMPLETED_TASKS == TASK_COUNT {
					close(failChan)
					return
				}

			} else {
				failChan <- task
			}

			continue
		}

		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
		fmt.Printf("\033[32m | Worker %d completed task %d\033[0m\n", id, task.id)

		// close the failed task channel if all the tasks have been completed
		COMPLETED_TASKS += 1
		if COMPLETED_TASKS == TASK_COUNT {
			close(failChan)
			return
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

	workerCount := 3

	// array of channels, one for each worker - used for round-robin task distribution
	taskChannels := make([]chan Task, workerCount)
	for i := 0; i < workerCount; i++ {
		taskChannels[i] = make(chan Task, len(tasks))
	}

	failChan := make(chan Task, len(tasks))
	failCounter := make(map[int]int)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i+1, taskChannels[i], failChan, &wg, failCounter)
	}

	// distribute the tasks using a round-robin approach
	currentWorkerID := 0
	for _, task := range tasks {
		taskChannels[currentWorkerID] <- task
		currentWorkerID = (currentWorkerID + 1) % workerCount // cycle worker id
	}

	// handler for failed tasks
	go func() {
		for failedTask := range failChan {
			fmt.Printf("\033[95m | Reassigning failed task %d to worker %d\033[0m\n", failedTask.id, currentWorkerID+1)
			taskChannels[currentWorkerID] <- failedTask
			currentWorkerID = (currentWorkerID + 1) % workerCount // cycle worker id
		}

		for _, taskChan := range taskChannels {
			close(taskChan)
		}
	}()

	wg.Wait()
	time.Sleep(time.Second)
	fmt.Print("\n\t\033[96;1;4mAll tasks completed\033[0m\n")
}
