package scheduler

import (
	"context"
	"log"
	"time"

	"betting-odds-scraper/internal/scraper"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron    *cron.Cron
	manager *scraper.Manager
}

func New(manager *scraper.Manager) *Scheduler {
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:    c,
		manager: manager,
	}
}

func (s *Scheduler) Start() {
	// Schedule scraping every 5 minutes
	_, err := s.cron.AddFunc("0 */5 * * * *", func() {
		log.Println("Starting scheduled scraping...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		results := s.manager.ScrapeAll(ctx)
		
		successCount := 0
		for _, result := range results {
			if result.Success {
				successCount++
			}
		}
		
		log.Printf("Scheduled scraping completed: %d/%d sites successful", successCount, len(results))
	})

	if err != nil {
		log.Printf("Failed to schedule scraping job: %v", err)
		return
	}

	// Schedule cleanup every hour
	_, err = s.cron.AddFunc("0 0 * * * *", func() {
		log.Println("Running cleanup tasks...")
		// Add cleanup logic here if needed
	})

	if err != nil {
		log.Printf("Failed to schedule cleanup job: %v", err)
	}

	s.cron.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) TriggerScrape() {
	go func() {
		log.Println("Manual scraping triggered...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		results := s.manager.ScrapeAll(ctx)
		
		successCount := 0
		for _, result := range results {
			if result.Success {
				successCount++
			}
		}
		
		log.Printf("Manual scraping completed: %d/%d sites successful", successCount, len(results))
	}()
}