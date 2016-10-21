package job

import (
	"errors"
	"log"
)

var (
	errTaskAlreadyExists = errors.New("Task already exists")
	errTaskNotExists     = errors.New("Task not exists")
)

// Scheduler interface represent a scheduler
type Scheduler interface {
	ScheduleTask(Task) error
	ScheduleJob(name string, job Job, cron ...string) error
	Stop(string) error
	All() []Task
	AllTask() []string
}

// scheduler is where you schedule all your jobs
type scheduler struct {
	tasks map[string]Task
}

// ScheduleJob register a task
func (s *scheduler) ScheduleTask(t Task) error {
	if _, ok := s.tasks[t.name()]; ok {
		return errTaskAlreadyExists
	}

	s.tasks[t.name()] = t
	// t.Start()
	t.run()
	return nil
}

// StopTask stop a task
func (s *scheduler) StopTask(t Task) error {
	if _, ok := s.tasks[t.name()]; !ok {
		return errTaskNotExists
	}
	delete(s.tasks, t.name())
	err := t.stop()
	t = nil
	return err
}

// ScheduleJob register a job
func (s *scheduler) ScheduleJob(name string, job Job, cron ...string) error {
	task, err := NewTask(name, job, cron...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return s.ScheduleTask(task)
}

// Stop will stop all registered jobs
func (s *scheduler) Stop(name string) error {
	if t, ok := s.tasks[name]; ok {
		return s.StopTask(t)
	}
	return errTaskNotExists
}

func (s *scheduler) All() []Task {
	var res []Task
	for _, v := range s.tasks {
		res = append(res, v)
	}
	return res
}

func (s *scheduler) AllTask() []string {
	var res []string
	for k := range s.tasks {
		res = append(res, k)
	}
	return res
}

// NewScheduler create a new scheduler
func NewScheduler() Scheduler {
	s := new(scheduler)
	s.tasks = make(map[string]Task)
	return s
}
