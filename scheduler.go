// Author: Pramesti Hatta K. @opam22
// 2022
package scheduler

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

type Unit int

const (
	Unknown Unit = iota
	Seconds
	Minutes
	Hour
	Day
)

type Writer interface {
	AddJob(job Job) error
	RemoveJob(id int) error
	Start()
	Close()
}

type Reader interface {
	GetJobs() []Job
	GetJob(id int) Job
}

type Schedule interface {
	Writer
	Reader
}

type Scheduler struct {
	Ctx     context.Context
	logger  Logger
	cleanup []func() error
	Jobs    []Job
	nextID  int
}

type Job struct {
	ID           int
	Name         string
	Every        int32
	Unit         Unit
	Task         func()
	RegisteredAt time.Time
}

// New() will return new instance of the scheduler
func New() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Scheduler{
		Ctx:    ctx,
		logger: DefaultLogger,
		nextID: 1,
	}

	s.cleanup = append(s.cleanup,
		func() error {
			cancel()
			s.logger.Info("cleaning up... cancel context")
			return nil
		},
		func() error {
			s.Jobs = []Job{}
			s.logger.Info("cleaning up... flushed all jobs")
			return nil
		})

	return s
}

// AddJob() can be used to add new job
// that we want to add as part of the scheduler
func (s *Scheduler) AddJob(job Job) error {
	if job.Name == "" {
		return errors.New("name is mandatory")
	}
	if job.Every == 0 {
		return errors.New("every is mandatory")
	}
	if job.Unit == Unknown {
		return errors.New("unit is mandatory")
	}

	j := Job{
		ID:           s.nextID,
		Name:         job.Name,
		Every:        job.Every,
		Unit:         job.Unit,
		Task:         job.Task,
		RegisteredAt: time.Now().Local(),
	}

	// convert every to second accross all unit
	// so we can get the exact second when this job will be triggered
	every := job.Every
	if job.Unit == Minutes {
		every = every * 60
	} else if job.Unit == Hour {
		every = every * 3600
	} else if job.Unit == Day {
		every = every * 86400
	}

	j.Every = every
	s.Jobs = append(s.Jobs, j)
	s.nextID = s.nextID + 1
	s.logger.Info("new job for scheduler was added")

	return nil
}

// GetJobs() can be used to get all registered jobs
func (s *Scheduler) GetJobs() []Job {
	return s.Jobs
}

// GetJob() can be used to get a job based on its id
func (s *Scheduler) GetJob(id int) Job {
	for _, job := range s.Jobs {
		if job.ID == id {
			return job
		}
	}

	return Job{}
}

// RemoveJob() can be used to remove a registered job based on id
func (s *Scheduler) RemoveJob(id int) error {
	for i, job := range s.Jobs {
		job := job
		if job.ID == id {
			s.Jobs = append(s.Jobs[:i], s.Jobs[i+1:]...)
		}
	}

	return nil
}

// Start() will kick off the scheduler and proceed all registered job
func (s *Scheduler) Start() {
	s.logger.Info("scheduler started")
	ticker := time.NewTicker(1 * time.Second)
	for ; true; <-ticker.C {
		if len(s.Jobs) == 0 {
			// no job, skip
			continue
		}

		wg := sync.WaitGroup{}
		for _, job := range s.Jobs {
			wg.Add(1)
			job := job
			go func() {
				defer wg.Done()
				s.run(job)
			}()
		}
		wg.Wait()
	}
}

// run() will run the job if its match the schedule
func (s *Scheduler) run(job Job) {
	if shouldRun(job) {
		job.Task()
	}
}

// shouldRun() detect if this is the time to run the job
func shouldRun(job Job) bool {
	currentTime := time.Now().Local()
	delta := currentTime.Sub(job.RegisteredAt)

	return int32(math.Ceil(delta.Seconds()))%job.Every == 0
}

// Close() will stop the scheduler
// and clean up the scheduler
func (s *Scheduler) Close() {
	for _, f := range s.cleanup {
		if err := f(); err != nil {
			s.logger.Error(err, "error when cleaning up the scheduler")
		}
	}
	s.cleanup = []func() error{}
}
