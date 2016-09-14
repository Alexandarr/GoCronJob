package job

import (
	"errors"
	"time"

	"github.com/robfig/revel/cache"
)

var (
	errJobAlreadyExists = errors.New("Job already exists")
)

// Scheduler interface represent a scheduler
type Scheduler interface {
	ScheduleJob(Jober, time.Duration) error
	Stop()
}

// scheduler is where you schedule all your jobs
type scheduler struct {
	store cache.Cache
	jobs  map[string]Jober
}

// ScheduleJob register a job
func (s *scheduler) ScheduleJob(j Jober, d time.Duration) error {
	if _, ok := s.jobs[j.Name()]; ok {
		return errJobAlreadyExists
	}
	s.jobs[j.Name()] = j
	return nil
}

// Stop will stop all registered jobs
func (s *scheduler) Stop() {
	for _, job := range s.jobs {
		job.Stop()
	}
}

// NewScheduler create a new scheduler
func NewScheduler(cache cache.Cache) Scheduler {
	s := new(scheduler)
	s.store = cache
	s.jobs = make(map[string]Jober)
	return s
}
