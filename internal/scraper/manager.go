package scraper

import (
	"context"
	"log"
	"sync"
	"time"

	"betting-odds-scraper/internal/config"
	"betting-odds-scraper/internal/models"
)

type Manager struct {
	config   *config.Config
	scrapers map[string]Scraper
	results  map[string][]models.ScrapeResult
	odds     map[string][]models.Odds
	matches  map[string]models.Match
	mutex    sync.RWMutex
}

type Scraper interface {
	GetSiteInfo() models.BettingSite
	ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error)
}

func NewManager(cfg *config.Config) *Manager {
	manager := &Manager{
		config:   cfg,
		scrapers: make(map[string]Scraper),
		results:  make(map[string][]models.ScrapeResult),
		odds:     make(map[string][]models.Odds),
		matches:  make(map[string]models.Match),
	}

	// Check if we should use demo mode (faster, no Chrome required)
	if cfg.LogLevel == "demo" {
		// Register demo scrapers for testing
		manager.RegisterScraper(NewDemoScraper("betika", "Betika"))
		manager.RegisterScraper(NewDemoScraper("sportpesa", "SportPesa"))
		manager.RegisterScraper(NewDemoScraper("betway", "Betway"))
		manager.RegisterScraper(NewDemoScraper("odibets", "Odibets"))
	} else {
		// Register real Kenyan betting site scrapers
		manager.RegisterScraper(NewBetikaScraper())
		manager.RegisterScraper(NewSportPesaScraper())
		manager.RegisterScraper(NewBetwayScraper())
		manager.RegisterScraper(NewOdibetsScraper())
	}

	return manager
}

func (m *Manager) RegisterScraper(scraper Scraper) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	siteInfo := scraper.GetSiteInfo()
	m.scrapers[siteInfo.ID] = scraper
	log.Printf("Registered scraper for %s", siteInfo.Name)
}

func (m *Manager) ScrapeAll(ctx context.Context) map[string]models.ScrapeResult {
	results := make(map[string]models.ScrapeResult)
	var wg sync.WaitGroup
	resultsChan := make(chan models.ScrapeResult, len(m.scrapers))

	// Limit concurrent scrapers
	semaphore := make(chan struct{}, m.config.MaxConcurrentScrapers)

	for siteID, scraper := range m.scrapers {
		wg.Add(1)
		go func(id string, s Scraper) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			result := m.scrapeWithTimeout(ctx, id, s)
			resultsChan <- result
		}(siteID, scraper)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		results[result.SiteID] = result
	}

	m.mutex.Lock()
	for siteID, result := range results {
		if m.results[siteID] == nil {
			m.results[siteID] = make([]models.ScrapeResult, 0)
		}
		m.results[siteID] = append(m.results[siteID], result)
		
		// Keep only last 10 results per site
		if len(m.results[siteID]) > 10 {
			m.results[siteID] = m.results[siteID][1:]
		}
	}
	m.mutex.Unlock()

	return results
}

func (m *Manager) scrapeWithTimeout(ctx context.Context, siteID string, scraper Scraper) models.ScrapeResult {
	start := time.Now()
	
	timeoutCtx, cancel := context.WithTimeout(ctx, m.config.RequestTimeout)
	defer cancel()

	matches, odds, err := scraper.ScrapeOdds(timeoutCtx)
	
	result := models.ScrapeResult{
		SiteID:    siteID,
		Success:   err == nil,
		Duration:  time.Since(start),
		ScrapedAt: time.Now(),
	}

	if err != nil {
		result.Error = err.Error()
		log.Printf("Failed to scrape %s: %v", siteID, err)
	} else {
		result.MatchCount = len(matches)
		result.OddsCount = len(odds)
		
		// Store results
		m.mutex.Lock()
		for _, match := range matches {
			m.matches[match.ID] = match
		}
		m.odds[siteID] = odds
		m.mutex.Unlock()
		
		log.Printf("Successfully scraped %s: %d matches, %d odds", siteID, len(matches), len(odds))
	}

	return result
}

func (m *Manager) GetBestOdds() []models.BestOdds {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	bestOddsMap := make(map[string]*models.BestOdds)

	// Process all odds for each match
	for _, odds := range m.odds {
		for _, odd := range odds {
			match, exists := m.matches[odd.MatchID]
			if !exists {
				continue
			}

			if bestOddsMap[odd.MatchID] == nil {
				bestOddsMap[odd.MatchID] = &models.BestOdds{
					Match:     match,
					AllOdds:   make([]models.Odds, 0),
					UpdatedAt: time.Now(),
				}
			}

			bestOdd := bestOddsMap[odd.MatchID]
			bestOdd.AllOdds = append(bestOdd.AllOdds, odd)

			// Update best odds
			if bestOdd.BestHomeWin == nil || odd.HomeWin > bestOdd.BestHomeWin.Value {
				bestOdd.BestHomeWin = &models.OddsComparison{
					Value:    odd.HomeWin,
					SiteID:   odd.SiteID,
					SiteName: odd.SiteName,
				}
			}

			if odd.Draw > 0 && (bestOdd.BestDraw == nil || odd.Draw > bestOdd.BestDraw.Value) {
				bestOdd.BestDraw = &models.OddsComparison{
					Value:    odd.Draw,
					SiteID:   odd.SiteID,
					SiteName: odd.SiteName,
				}
			}

			if bestOdd.BestAwayWin == nil || odd.AwayWin > bestOdd.BestAwayWin.Value {
				bestOdd.BestAwayWin = &models.OddsComparison{
					Value:    odd.AwayWin,
					SiteID:   odd.SiteID,
					SiteName: odd.SiteName,
				}
			}
		}
	}

	result := make([]models.BestOdds, 0, len(bestOddsMap))
	for _, bestOdd := range bestOddsMap {
		result = append(result, *bestOdd)
	}

	return result
}

func (m *Manager) GetScrapeResults() map[string][]models.ScrapeResult {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	results := make(map[string][]models.ScrapeResult)
	for k, v := range m.results {
		results[k] = make([]models.ScrapeResult, len(v))
		copy(results[k], v)
	}
	return results
}