package job

import (
	"log"
	"time"

	"github.com/f2prateek/clock"
	"github.com/gorhill/cronexpr"
)

var c = clock.Default()

type ticker struct {
	c       chan time.Time // The channel on which the ticks are delivered.
	doneC   chan struct{}  // The channel on which the stop signal is delivered.
	pauseC  chan struct{}  // The channel on which the pause signal is delivered.
	resumeC chan struct{}  // The channel on which the resume signal is delivered.
	expr    *cronexpr.Expression
}

// Parse returns a new Ticker containing a channel that will send the time
// with a period specified by the spec argument. Stop the ticker to release
// associated resources.
func parse(spec string) (*ticker, error) {
	expr, err := cronexpr.Parse(spec)
	if err != nil {
		return nil, err
	}

	tickerC := make(chan time.Time, 1)
	ticker := &ticker{
		c:       tickerC,
		doneC:   make(chan struct{}, 1),
		pauseC:  make(chan struct{}, 1),
		resumeC: make(chan struct{}, 1),
		expr:    expr,
	}

	return ticker, nil
}

// Start start a ticker.
func (t *ticker) start() error {
	go func() {
		for {
			next := t.expr.Next(c.Now())
			log.Println(next)
			select {
			case <-time.After(next.Sub(c.Now())):
				t.c <- c.Now()
			case <-t.pauseC:
				select {
				case <-t.resumeC:
					continue
				case <-t.doneC:
					return
				}
			case <-t.doneC:
				break
			}
		}
	}()
	return nil
}

// Pause pause a ticker. After pause, no more ticks will be sent.
func (t *ticker) pause() error {
	t.pauseC <- struct{}{}
	return nil
}

// Resume Resume a paused ticker.
func (t *ticker) resume() error {
	t.resumeC <- struct{}{}
	return nil
}

// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does
// not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t *ticker) stop() error {
	t.doneC <- struct{}{}
	return nil
}
