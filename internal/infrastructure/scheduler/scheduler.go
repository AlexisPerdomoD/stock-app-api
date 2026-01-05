package scheduler

import (
	"context"
	"fmt"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type Scheduler struct {
	intance *cron.Cron
	jobs    map[string]cron.EntryID
}

/*
interval default is 24h
*/
func (sc *Scheduler) AddStockSourceService(
	serviceName string,
	uc *usecases.RegisterStocks,
	timeout time.Duration,
	intervalDuration *time.Duration,
) {

	if uc == nil {
		log.Fatalln("bad impl: scheduler required args was passed nil for AddStockSourceService")
	}

	if sc.jobs[serviceName] != 0 {
		log.Println("[CRON] already added stock sourcing at:", sc.jobs[serviceName], "source:", serviceName, "skipping")
		return
	}

	interval := "@every 24h"
	if intervalDuration != nil {
		interval = fmt.Sprintf("@every %s", intervalDuration.String())
	}

	id, err := sc.intance.AddFunc(interval, func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		log.Println("[CRON] starting stock source service ", serviceName)

		limitDate := time.Now().AddDate(0, 0, 1)
		if intervalDuration != nil {
			limitDate = time.Now().Add(*intervalDuration)
		}

		inserts, err := uc.Execute(ctx, serviceName, &limitDate)
		if err != nil {
			log.Println("[CRON] error executing stock source service ", serviceName)
			log.Println("[CRON] error: ", err)
			return
		}

		log.Println("[CRON] inserted rows:", inserts)

		delete(sc.jobs, serviceName)
		log.Println("[CRON] finished stock source service:", serviceName)
	})

	if err != nil {
		log.Println(err)
		log.Fatalln("bad impl: scheduler failed to add func ", interval)
	}

	sc.jobs[serviceName] = id
	log.Println("[CRON] added stock sourcing at:", interval, "source:", serviceName, "with id:", id)
}

func (sc *Scheduler) StartOnBackground() {
	sc.intance.Start()
	log.Println("[CRON] started")
}

func New() *Scheduler {
	return &Scheduler{
		intance: cron.New(),
		jobs:    make(map[string]cron.EntryID),
	}
}
