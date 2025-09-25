package models

import (
	"time"
)

// BettingSite represents a betting platform
type BettingSite struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Active   bool   `json:"active"`
	LastScrape time.Time `json:"last_scrape"`
}

// Match represents a sports match
type Match struct {
	ID          string    `json:"id"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	Sport       string    `json:"sport"`
	League      string    `json:"league"`
	MatchTime   time.Time `json:"match_time"`
	Status      string    `json:"status"`
}

// Odds represents betting odds for a match
type Odds struct {
	ID         string    `json:"id"`
	MatchID    string    `json:"match_id"`
	SiteID     string    `json:"site_id"`
	SiteName   string    `json:"site_name"`
	HomeWin    float64   `json:"home_win"`
	Draw       float64   `json:"draw,omitempty"`
	AwayWin    float64   `json:"away_win"`
	Over25     float64   `json:"over_2_5,omitempty"`
	Under25    float64   `json:"under_2_5,omitempty"`
	BTTS       float64   `json:"btts,omitempty"`
	ScrapedAt  time.Time `json:"scraped_at"`
}

// BestOdds represents the best odds found across all sites
type BestOdds struct {
	Match       Match              `json:"match"`
	BestHomeWin *OddsComparison    `json:"best_home_win"`
	BestDraw    *OddsComparison    `json:"best_draw,omitempty"`
	BestAwayWin *OddsComparison    `json:"best_away_win"`
	BestOver25  *OddsComparison    `json:"best_over_2_5,omitempty"`
	BestUnder25 *OddsComparison    `json:"best_under_2_5,omitempty"`
	BestBTTS    *OddsComparison    `json:"best_btts,omitempty"`
	AllOdds     []Odds             `json:"all_odds"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// OddsComparison represents the best odds for a specific market
type OddsComparison struct {
	Value    float64 `json:"value"`
	SiteID   string  `json:"site_id"`
	SiteName string  `json:"site_name"`
}

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	SiteID    string    `json:"site_id"`
	Success   bool      `json:"success"`
	MatchCount int      `json:"match_count"`
	OddsCount int       `json:"odds_count"`
	Error     string    `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	ScrapedAt time.Time `json:"scraped_at"`
}