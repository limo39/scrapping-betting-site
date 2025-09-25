package scraper

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"betting-odds-scraper/internal/models"
)

// DemoScraper provides sample data without actually scraping websites
type DemoScraper struct {
	siteInfo models.BettingSite
}

func NewDemoScraper(siteID, siteName string) *DemoScraper {
	return &DemoScraper{
		siteInfo: models.BettingSite{
			ID:     siteID,
			Name:   siteName,
			URL:    fmt.Sprintf("https://www.%s.com", siteID),
			Active: true,
		},
	}
}

func (d *DemoScraper) GetSiteInfo() models.BettingSite {
	return d.siteInfo
}

func (d *DemoScraper) ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error) {
	var matches []models.Match
	var odds []models.Odds

	// Simulate some processing time
	time.Sleep(time.Duration(rand.Intn(2000)+500) * time.Millisecond)

	// Sample matches with realistic team names
	sampleMatches := []struct {
		home, away string
		league     string
	}{
		{"Arsenal", "Chelsea", "Premier League"},
		{"Manchester United", "Liverpool", "Premier League"},
		{"Barcelona", "Real Madrid", "La Liga"},
		{"Bayern Munich", "Borussia Dortmund", "Bundesliga"},
		{"PSG", "Marseille", "Ligue 1"},
		{"Juventus", "AC Milan", "Serie A"},
		{"Tottenham", "Manchester City", "Premier League"},
		{"Atletico Madrid", "Sevilla", "La Liga"},
		{"Inter Milan", "Napoli", "Serie A"},
		{"Leicester City", "West Ham", "Premier League"},
		{"Valencia", "Villarreal", "La Liga"},
		{"RB Leipzig", "Bayer Leverkusen", "Bundesliga"},
	}

	// Generate random matches (3-8 matches per site)
	numMatches := rand.Intn(6) + 3
	selectedMatches := make([]int, numMatches)
	
	// Select random matches without duplicates
	used := make(map[int]bool)
	for i := 0; i < numMatches; i++ {
		for {
			idx := rand.Intn(len(sampleMatches))
			if !used[idx] {
				selectedMatches[i] = idx
				used[idx] = true
				break
			}
		}
	}

	for i, matchIdx := range selectedMatches {
		sample := sampleMatches[matchIdx]
		
		matchID := fmt.Sprintf("%s_%s_vs_%s_%d", 
			d.siteInfo.ID,
			strings.ReplaceAll(strings.ToLower(sample.home), " ", "_"),
			strings.ReplaceAll(strings.ToLower(sample.away), " ", "_"),
			time.Now().Unix()+int64(i))

		match := models.Match{
			ID:        matchID,
			HomeTeam:  sample.home,
			AwayTeam:  sample.away,
			Sport:     "football",
			League:    sample.league,
			MatchTime: time.Now().Add(time.Duration(24+i*6) * time.Hour),
			Status:    "upcoming",
		}

		matches = append(matches, match)

		// Generate realistic odds with some variation per site
		baseHomeOdds := 1.5 + rand.Float64()*2.0  // 1.5 to 3.5
		baseDraw := 3.0 + rand.Float64()*1.0      // 3.0 to 4.0
		baseAwayOdds := 1.5 + rand.Float64()*2.0  // 1.5 to 3.5

		// Add site-specific variation
		siteVariation := rand.Float64()*0.2 - 0.1 // -0.1 to +0.1
		
		homeOdds := baseHomeOdds + siteVariation
		drawOdds := baseDraw + siteVariation
		awayOdds := baseAwayOdds + siteVariation

		// Ensure minimum odds
		if homeOdds < 1.1 { homeOdds = 1.1 }
		if drawOdds < 2.5 { drawOdds = 2.5 }
		if awayOdds < 1.1 { awayOdds = 1.1 }

		odd := models.Odds{
			ID:        fmt.Sprintf("%s_odds", matchID),
			MatchID:   matchID,
			SiteID:    d.siteInfo.ID,
			SiteName:  d.siteInfo.Name,
			HomeWin:   homeOdds,
			Draw:      drawOdds,
			AwayWin:   awayOdds,
			ScrapedAt: time.Now(),
		}
		odds = append(odds, odd)
	}

	return matches, odds, nil
}