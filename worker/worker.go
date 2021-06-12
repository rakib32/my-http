package worker

import (
	"fmt"
	"sync"
)

type Handler func(interface{}) (interface{}, error)

type Job struct {
	Payload interface{}
	Handler Handler
}

func (j *Job) Handle() (interface{}, error) {
	resp, err := j.Handler(j.Payload)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

type Worker struct {
	ID         int
	WorkerPool chan chan Job
	Job        chan Job
	Quit       chan int
}

func NewWorker(id int, workerPool chan chan Job, quit chan int) Worker {
	return Worker{
		ID:         id,
		WorkerPool: workerPool,
		Job:        make(chan Job),
		Quit:       quit,
	}
}
func startWorker(w Worker) {
	for {
		// Add job channel to the worker queue.
		w.WorkerPool <- w.Job
		select {
		case job := <-w.Job:
			if _, err := job.Handle(); err != nil {
				fmt.Printf("Error %v", err)
			}
		case <-w.Quit:
			fmt.Printf("quitting %v \n", w.ID)
			return
		}
	}
}

func (w Worker) Start(wg *sync.WaitGroup) {
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		startWorker(w)
	}(wg)
}
