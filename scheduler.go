package scheduler

import (
	"math"
	"time"
)

type Unit int

const (
	Seconds Unit = iota
	Minutes
	Hour
	Day
)

type Handler struct {
	Jobs []Job
}

type Job struct {
	Name         string
	Every        int32
	Unit         Unit
	Task         func()
	RegisteredAt time.Time
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) AddJob(job Job) {
	j := Job{
		Name:         job.Name,
		Every:        job.Every,
		Unit:         job.Unit,
		Task:         job.Task,
		RegisteredAt: time.Now(),
	}

	h.Jobs = append(h.Jobs, j)
}

func (h *Handler) Start() {
	for range time.Tick(time.Second * 1) {
		for _, j := range h.Jobs {
			h.Run(j)
		}
	}
}

func (h *Handler) Run(job Job) {
	if ShouldRun(job) {
		job.Task()
	}
}

func ShouldRun(job Job) bool {
	currentTime := time.Now()
	delta := currentTime.Sub(job.RegisteredAt)

	if job.Unit == Seconds {
		if int32(math.Ceil(delta.Seconds()))%job.Every == 0 {
			return true
		}
	} else if job.Unit == Minutes {
		if int32(math.Ceil(delta.Seconds()))%(job.Every*60) == 0 {
			return true
		}
	} else if job.Unit == Hour {
		if int32(math.Ceil(delta.Seconds()))%(job.Every*3600) == 0 {
			return true
		}
	}

	return false
}
