package mini

// Kafka runner for the scheduler

import (
	"fmt"
	"strconv"
	"time"

	hmapcoll "github.com/MingyangXiaTim/kafka-scheduler/internal/collector/hmap"
	"github.com/MingyangXiaTim/kafka-scheduler/internal/test"
	"github.com/MingyangXiaTim/kafka-scheduler/schedule/kafka"
	log "github.com/sirupsen/logrus"

	"github.com/MingyangXiaTim/kafka-scheduler/apiserver/rest"
	"github.com/MingyangXiaTim/kafka-scheduler/config"
	"github.com/MingyangXiaTim/kafka-scheduler/internal/store/hmap"
	"github.com/MingyangXiaTim/kafka-scheduler/scheduler"
)

const (
	ID1 = 1
	ID2 = 2
	ID3 = 3
)

func GetSchedule(id int) kafka.Schedule {
	scheduleID := fmt.Sprintf("schedule-%d", id)
	targetID := strconv.Itoa(id)
	value := fmt.Sprintf("value %d", id)

	epoch := time.Now().Add(time.Duration(id) * time.Minute).Unix()
	return kafka.Schedule{Message: test.FullMessage("schedules", scheduleID, value, epoch, "topic", targetID)}
}

type Runner struct {
	stopChan chan bool
}

func NewRunner() *Runner {
	return &Runner{
		stopChan: make(chan bool),
	}
}

func (r Runner) Close() error {
	r.stopChan <- true
	return nil
}

func (r *Runner) Start() error {
	store := hmap.New()

	store.Add(GetSchedule(ID1))
	store.Add(GetSchedule(ID2))
	store.Add(GetSchedule(ID3))

	handler := NewHandler(store)

	sch := scheduler.New(store, hmapcoll.New(), nil)
	sch.Start(scheduler.StartOfToday())

	srv := rest.New(&sch)
	srv.Start(config.ServerAddr())

	events := sch.Events()

loop:
	for {
		select {
		case event, open := <-events:
			if !open {
				break loop
			}
			handler.Handle(event)
		case <-r.stopChan:
			err := srv.Stop()
			if err != nil {
				log.Errorf("error when stopping api server: %v", err)
			}
			sch.Close()
		}
	}

	return nil
}
