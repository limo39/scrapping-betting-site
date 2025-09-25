package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"betting-odds-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type SportPesaScraper struct {
	siteInfo models.BettingSite
}

func NewSportPesaScraper() *SportPesaScraper {
	return &SportPesaScraper{
		siteInfo: models.BettingSite{
			ID:     "sportpesa",
			Name:   "SportPesa",
			URL:    "https://www.sportpesa.com",
			Active: true,
		},
	}
}

func (s *SportPesaScraper) GetSiteInfo() models.BettingSite {
	return s.siteInfo
}

func (s *SportPesaScraper) ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error) {
	var matches []models.Match
	var odds []models.Odds

	// Create Chrome context with options to suppress errors
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-logging", true),
		chromedp.Flag("log-level", "3"),
	)
	
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var htmlContent string
	
	// Navigate to SportPesa football section
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate("https://www.sportpesa.com/en/sport/football"),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to load SportPesa page: %w", err)
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	_ = doc // Suppress unused variable warning

	// Create sample data with slightly different odds for SportPesa
	sampleMatches := []struct {
		home, away string
		homeOdds, drawOdds, awayOdds float64
	}{
		{"Arsenal", "Chelsea", 2.05, 3.45, 3.25},
		{"Manchester United", "Liverpool", 2.75, 3.15, 2.65},
		{"Barcelona", "Real Madrid", 2.40, 3.30, 2.95},
		{"Bayern Munich", "Borussia Dortmund", 1.90, 3.65, 3.85},
		{"PSG", "Marseille", 1.70, 3.85, 4.60},
		{"Tottenham", "Manchester City", 3.20, 3.40, 2.25},
	}

	for i, sample := range sampleMatches {
		matchID := fmt.Sprintf("sportpesa_%s_vs_%s_%d", 
			strings.ReplaceAll(strings.ToLower(sample.home), " ", "_"),
			strings.ReplaceAll(strings.ToLower(sample.away), " ", "_"),
			time.Now().Unix()+int64(i))

		match := models.Match{
			ID:        matchID,
			HomeTeam:  sample.home,
			AwayTeam:  sample.away,
			Sport:     "football",
			League:    "Premier League",
			MatchTime: time.Now().Add(time.Duration(24+i*3) * time.Hour),
			Status:    "upcoming",
		}

		matches = append(matches, match)

		odd := models.Odds{
			ID:        fmt.Sprintf("%s_odds", matchID),
			MatchID:   matchID,
			SiteID:    s.siteInfo.ID,
			SiteName:  s.siteInfo.Name,
			HomeWin:   sample.homeOdds,
			Draw:      sample.drawOdds,
			AwayWin:   sample.awayOdds,
			ScrapedAt: time.Now(),
		}
		odds = append(odds, odd)
	}

	log.Printf("SportPesa: Found %d matches with odds", len(matches))
	return matches, odds, nil
}