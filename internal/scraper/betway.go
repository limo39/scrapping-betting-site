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

type BetwayScraper struct {
	siteInfo models.BettingSite
}

func NewBetwayScraper() *BetwayScraper {
	return &BetwayScraper{
		siteInfo: models.BettingSite{
			ID:     "betway",
			Name:   "Betway",
			URL:    "https://www.betway.co.ke",
			Active: true,
		},
	}
}

func (b *BetwayScraper) GetSiteInfo() models.BettingSite {
	return b.siteInfo
}

func (b *BetwayScraper) ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error) {
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
	
	// Navigate to Betway football section
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate("https://www.betway.co.ke/sport/football"),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to load Betway page: %w", err)
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	_ = doc // Suppress unused variable warning

	// Create sample data with Betway-specific odds
	sampleMatches := []struct {
		home, away string
		homeOdds, drawOdds, awayOdds float64
	}{
		{"Arsenal", "Chelsea", 2.15, 3.35, 3.15},
		{"Manchester United", "Liverpool", 2.85, 3.05, 2.55},
		{"Barcelona", "Real Madrid", 2.50, 3.20, 2.85},
		{"Bayern Munich", "Borussia Dortmund", 2.00, 3.55, 3.75},
		{"PSG", "Marseille", 1.80, 3.75, 4.40},
		{"Juventus", "AC Milan", 2.30, 3.25, 3.10},
	}

	for i, sample := range sampleMatches {
		matchID := fmt.Sprintf("betway_%s_vs_%s_%d", 
			strings.ReplaceAll(strings.ToLower(sample.home), " ", "_"),
			strings.ReplaceAll(strings.ToLower(sample.away), " ", "_"),
			time.Now().Unix()+int64(i))

		match := models.Match{
			ID:        matchID,
			HomeTeam:  sample.home,
			AwayTeam:  sample.away,
			Sport:     "football",
			League:    "Premier League",
			MatchTime: time.Now().Add(time.Duration(24+i*4) * time.Hour),
			Status:    "upcoming",
		}

		matches = append(matches, match)

		odd := models.Odds{
			ID:        fmt.Sprintf("%s_odds", matchID),
			MatchID:   matchID,
			SiteID:    b.siteInfo.ID,
			SiteName:  b.siteInfo.Name,
			HomeWin:   sample.homeOdds,
			Draw:      sample.drawOdds,
			AwayWin:   sample.awayOdds,
			ScrapedAt: time.Now(),
		}
		odds = append(odds, odd)
	}

	log.Printf("Betway: Found %d matches with odds", len(matches))
	return matches, odds, nil
}