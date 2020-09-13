package scheduler

import (
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type Scheduler struct {
	m       sync.Mutex
	wg      sync.WaitGroup
	started bool
	stopped bool

	precision time.Duration
	schedules []*Schedule
	FanIn     chan int
	stopChan  chan struct{}
}

func (s *Scheduler) Add(schedule *Schedule) {
	if schedule == nil {
		return
	}

	s.m.Lock()
	s.schedules = append(s.schedules, schedule)
	s.m.Unlock()
}

func (s *Scheduler) start() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic: it may happen at time of stopping the scheduler.\nStackTrace: ", string(debug.Stack()))
		}

		s.wg.Done()
	}()

	s.m.Lock()
	s.started = true
	s.wg.Add(1)
	s.m.Unlock()

	keepRunning := true

	go func() {
		select {
		case <-s.stopChan:
			keepRunning = false
		}
	}()

	for keepRunning {
		i := 0
		for i < len(s.schedules) {
			schedule := s.schedules[i]

			select {
			case <-schedule.sig:
				s.FanIn <- schedule.ID

				s.m.Lock()
				s.schedules = append(s.schedules[:i], s.schedules[i+1:]...)
				s.m.Unlock()
			default:
				i += 1
			}
		}
		time.Sleep(s.precision)
	}
}

func (s *Scheduler) Start() {
	if s.started {
		log.Println("Scheduler already started.")
		return
	}

	go s.start()
}

func (s *Scheduler) Stop() {
	if s.stopped {
		log.Println("Scheduler already stopped.")
		return
	}

	s.stopChan <- struct{}{}

	s.m.Lock()
	schedules := s.schedules
	s.schedules = nil
	close(s.FanIn)
	s.m.Unlock()

	for _, s := range schedules {
		s.timer.Stop()
	}

	s.wg.Wait()
}

type SchedulerBuilder struct {
	scheduler *Scheduler
}

func NewSchedulerBuilder() *SchedulerBuilder {
	return &SchedulerBuilder{
		scheduler: &Scheduler{
			schedules: make([]*Schedule, 0),
			FanIn:     make(chan int),
			stopChan:  make(chan struct{}),
		},
	}
}

func (s *SchedulerBuilder) WithPrecision(duration time.Duration) *SchedulerBuilder {
	s.scheduler.precision = duration
	return s
}

func (s *SchedulerBuilder) Build() *Scheduler {
	return s.scheduler
}
