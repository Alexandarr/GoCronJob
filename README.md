# GoCronJob

GoCronJob is a job scheduler support cron expression

## Immediately Job

```go
import (
    . github.com/philchia/GoCronJob
)

func Example() {
	j := func() {
		log.Println("Hello, world!")
	}

	s := NewScheduler()
	s.ScheduleJob("Hello", j)
}
```

## Cron Job

```go
import (
    . github.com/philchia/GoCronJob
)

func Example() {
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
```