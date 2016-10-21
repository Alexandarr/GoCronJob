package job

import (
	"errors"
	"fmt"
	"log"
)

// Job is a alias name of a runable function
type Job func()

// Status represent task status
type Status int

// Description will return the description of a status in string
func (s Status) Description() string {
	switch s {
	case Ready:
		return "ready"
	case Running:
		return "running"
	case Paused:
		return "paused"
	case Stoped:
		return "stoped"
	default:
		return "unknow"
	}
}

const (
	// Ready : 		Task is ready to start
	Ready Status = iota
	// Running : 	Task is running
	Running
	// Paused : 	Task is paused, may resume or stop
	Paused
	// Stoped : 	Task is stoped, will not response to any method
	Stoped
)

// Task represent a task
type Task interface {
	run() error
	stop() error
	pause() error
	resume() error
	status() Status
	name() string
	description() string
}

type task struct {
	taskName   string
	process    float32
	job        Job
	taskStatus Status

	doneC chan struct{}
	t     *ticker
}

// NewTask generate a new Task interface
func NewTask(name string, job Job, cron ...string) (Task, error) {
	var t *ticker
	var err error

	if len(cron) > 0 {
		t, err = parse(cron[0])
		if err != nil {
			return nil, err
		}
	}

	task := new(task)
	task.t = t
	task.taskName = name
	task.job = job
	log.Println("Make job")
	task.doneC = make(chan struct{})

	return task, nil
}

func (task *task) run() error {
	if task.t == nil {
		log.Println("Run task immediately")
		task.taskStatus = Running
		task.job()
		task.taskStatus = Stoped
		return nil
	}

	go func() {
		log.Println("Run task repeatly")

		for {
			select {
			case <-task.t.c:
				log.Println("Ticker")
				switch task.taskStatus {
				case Ready, Paused:
					func() {
						task.taskStatus = Running
						task.job()
						task.taskStatus = Ready
					}()
				case Running:
					continue
				default:
					return
				}
			case <-task.doneC:
				return
			}
		}
	}()
	return task.t.start()
}

func (task *task) stop() error {

	return task.t.stop()
}

func (task *task) pause() error {
	if task.t == nil {
		return errors.New("can not stop a noncron task")
	}
	return nil
}

func (task *task) resume() error {
	if task.t == nil {
		return errors.New("can not resume a noncron task")
	}
	return task.t.resume()
}

func (task *task) status() Status {
	return task.taskStatus
}

func (task *task) name() string {
	return task.taskName
}

func (task *task) description() string {
	return fmt.Sprintf(`\n***********************\nTask: %s\nStatus: %s\nProcess: %f\n***********************\n`, task.taskName, task.taskStatus.Description(), task.process)
}
