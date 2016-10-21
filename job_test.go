package job

import (
	"log"
	"testing"
	"time"
)

func TestImmediateJob(t *testing.T) {

	j := func() {
		log.Println("Hello, world!")
	}

	s := NewScheduler()
	s.ScheduleJob("Hello", j)
}

func TestRepeatJob(t *testing.T) {
	j := func() {
		log.Println("Hello, world!")
	}

	s := NewScheduler()
	err := s.ScheduleJob("Hello", j, "* * * * * * *")
	if err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
}
