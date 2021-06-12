package worker

import (
	"sync"
)

type Dispatcher struct {
	WorkerPool chan chan Job
	maxWorkers int
	JobQueue   chan Job
	Quit       chan int
	Wg         sync.WaitGroup
}

func NewDispatcher(maxWorkers int, jobQueue chan Job) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		JobQueue:   jobQueue,
		Quit:       make(chan int),
		Wg:         sync.WaitGroup{},
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		d.Wg.Add(1)
		worker := NewWorker(i, d.WorkerPool, d.Quit)
		worker.Start(&d.Wg)
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				worker := <-d.WorkerPool
				//fmt.Println("Job dispatched: ", job.Payload)

				// dispatch the job to the worker job channel
				worker <- job
			}(job)
		}
	}
}

func (d *Dispatcher) Stop() {
	for i := 0; i < d.maxWorkers; i++ {
		d.Quit <- 1
	}
	d.Wg.Wait()
	close(d.JobQueue)
	close(d.WorkerPool)
}

func (d *Dispatcher) AddJob(job Job) {
	d.JobQueue <- job
}

func Init(maxJobQueue, maxWorkerCount int) *Dispatcher {
	q := make(chan Job, maxJobQueue)
	dispatcher := NewDispatcher(maxWorkerCount, q)
	dispatcher.Run()
	return dispatcher
}
