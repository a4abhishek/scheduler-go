package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/a4abhishek/schedular-go"
)

type schedule struct {
	userID  int
	message string
	in      time.Duration
}

func main() {
	schedules := map[int]*schedule{}

	schedules[0] = &schedule{
		userID:  0,
		message: "Message 1",
		in:      5 * time.Second,
	}

	schedules[1] = &schedule{
		userID:  1,
		message: "Message 2",
		in:      1 * time.Second,
	}

	sdr := scheduler.NewSchedulerBuilder().
		WithPrecision(time.Second).
		Build()

	sdr.Start()

	sdr.Add(schedules[0].GetSchedule())
	sdr.Add(schedules[1].GetSchedule())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(7 * time.Second)
		sdr.Stop()
		fmt.Println("Stopped the scheduler.")
		wg.Done()
	}()

	for id := range sdr.FanIn {
		schedule := schedules[id]
		fmt.Printf("User-ID: %d, Message: %q\n", schedule.userID, schedule.message)
	}

	wg.Wait()
}

func (s *schedule) GetSchedule() *scheduler.Schedule {
	return scheduler.NewScheduleBuilder(s.userID).
		WithDuration(s.in).
		Build()
}
