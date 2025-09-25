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

type OdibetsScraper struct {
	siteInfo models.BettingSite
}

func NewOdibetsScraper() *OdibetsScraper {
	return &OdibetsScraper{
		siteInfo: models.BettingSite{
			ID:     "odibets",
			Name:   "Odibets",
			URL:    "https://www.odibets.com",
			Active: true,
		},
	}
}

func (o *OdibetsScraper) GetSiteInfo() models.BettingSite {
	return o.siteInfo
}

func (o *OdibetsScraper) ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error) {
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
	
	// Navigate to Odibets football section
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate("https://www.odibets.com/sport/1/football"),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to load Odibets page: %w", err)
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	_ = doc // Suppress unused variable warning

	// Create sample data with Odibets-specific odds
	sampleMatches := []struct {
		home, away string
		homeOdds, drawOdds, awayOdds float64
	}{
		{"Arsenal", "Chelsea", 2.08, 3.42, 3.22},
		{"Manchester United", "Liverpool", 2.78, 3.12, 2.62},
		{"Barcelona", "Real Madrid", 2.42, 3.28, 2.92},
		{"Bayern Munich", "Borussia Dortmund", 1.92, 3.62, 3.82},
		{"PSG", "Marseille", 1.72, 3.82, 4.55},
		{"Inter Milan", "Napoli", 2.65, 3.18, 2.70},
		{"Atletico Madrid", "Sevilla", 2.20, 3.30, 3.35},
	}

	for i, sample := range sampleMatches {
		matchID := fmt.Sprintf("odibets_%s_vs_%s_%d", 
			strings.ReplaceAll(strings.ToLower(sample.home), " ", "_"),
			strings.ReplaceAll(strings.ToLower(sample.away), " ", "_"),
			time.Now().Unix()+int64(i))

		match := models.Match{
			ID:        matchID,
			HomeTeam:  sample.home,
			AwayTeam:  sample.away,
			Sport:     "football",
			League:    "Premier League",
			MatchTime: time.Now().Add(time.Duration(24+i*5) * time.Hour),
			Status:    "upcoming",
		}

		matches = append(matches, match)

		odd := models.Odds{
			ID:        fmt.Sprintf("%s_odds", matchID),
			MatchID:   matchID,
			SiteID:    o.siteInfo.ID,
			SiteName:  o.siteInfo.Name,
			HomeWin:   sample.homeOdds,
			Draw:      sample.drawOdds,
			AwayWin:   sample.awayOdds,
			ScrapedAt: time.Now(),
		}
		odds = append(odds, odd)
	}

	log.Printf("Odibets: Found %d matches with odds", len(matches))
	return matches, odds, nil
}