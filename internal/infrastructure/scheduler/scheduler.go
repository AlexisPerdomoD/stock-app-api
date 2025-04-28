package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	cron "github.com/robfig/cron/v3"
)

type Scheduler struct {
	intance *cron.Cron
	jobs    map[string]cron.EntryID
}

/*
interval default is 24h
*/
func (sc *Scheduler) AddStockSourceService(
	s domain.SourceStockService,
	uc *usecase.RegisterStocksUseCase,
	timeout time.Duration,
	itv *time.Duration,
) {

	if s == nil || uc == nil {
		log.Fatalln("bad impl: scheduler required args was passed nil for AddStockSourceService")
	}

	if sc.jobs[s.Name()] != 0 {
		log.Println("[CRON] already added stock sourcing at:", sc.jobs[s.Name()], "source:", s.Name(), "skipping")
		return
	}

	interval := "@every 24h"
	limitDate := time.Now().AddDate(0, 0, 1)
	if itv != nil {
		interval = fmt.Sprintf("@every %s", itv.String())
		limitDate = time.Now().Add(-*itv)
	}

	id, err := sc.intance.AddFunc(interval, func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		log.Println("[CRON] starting stock source service ", s.Name())

		inserts, err := uc.Execute(ctx, s, &limitDate)
		if err != nil {
			log.Println("[CRON] error executing stock source service ", s.Name())
			log.Println("[CRON] error: ", err)
			return
		}

		log.Println("[CRON] inserted rows:", inserts)

		delete(sc.jobs, s.Name())
		log.Println("[CRON] finished stock source service:", s.Name())
	})

	if err != nil {
		log.Println(err)
		log.Fatalln("bad impl: scheduler failed to add func ", interval)
	}

	sc.jobs[s.Name()] = id
	log.Println("[CRON] added stock sourcing at:", interval, "source:", s.Name(), "with id:", id)
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
