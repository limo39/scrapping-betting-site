package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"betting-odds-scraper/internal/config"
	"betting-odds-scraper/internal/scraper"
)

func main() {
	fmt.Println("🏆 Testing Kenya Betting Odds Scraper")
	fmt.Println("=====================================")

	// Initialize configuration
	cfg := config.New()
	cfg.RequestTimeout = 60 * time.Second // Longer timeout for testing

	// Initialize scraper manager
	manager := scraper.NewManager(cfg)

	// Test scraping
	fmt.Println("\n🔄 Starting test scrape...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	results := manager.ScrapeAll(ctx)

	fmt.Println("\n📊 Scraping Results:")
	fmt.Println("====================")

	totalMatches := 0
	totalOdds := 0
	successfulSites := 0

	for siteID, result := range results {
		status := "❌ FAILED"
		if result.Success {
			status = "✅ SUCCESS"
			successfulSites++
		}

		fmt.Printf("%-12s %s - Matches: %d, Odds: %d, Duration: %v\n",
			siteID, status, result.MatchCount, result.OddsCount, result.Duration)

		if result.Error != "" {
			fmt.Printf("             Error: %s\n", result.Error)
		}

		totalMatches += result.MatchCount
		totalOdds += result.OddsCount
	}

	fmt.Println("\n📈 Summary:")
	fmt.Println("===========")
	fmt.Printf("Successful Sites: %d/%d\n", successfulSites, len(results))
	fmt.Printf("Total Matches: %d\n", totalMatches)
	fmt.Printf("Total Odds: %d\n", totalOdds)

	if totalOdds > 0 {
		fmt.Println("\n🏆 Best Odds Found:")
		fmt.Println("===================")

		bestOdds := manager.GetBestOdds()
		for i, match := range bestOdds {
			if i >= 5 { // Show only first 5 matches
				break
			}

			fmt.Printf("\n%s vs %s\n", match.Match.HomeTeam, match.Match.AwayTeam)
			
			if match.BestHomeWin != nil {
				fmt.Printf("  Home Win: %.2f (%s)\n", match.BestHomeWin.Value, match.BestHomeWin.SiteName)
			}
			
			if match.BestDraw != nil {
				fmt.Printf("  Draw:     %.2f (%s)\n", match.BestDraw.Value, match.BestDraw.SiteName)
			}
			
			if match.BestAwayWin != nil {
				fmt.Printf("  Away Win: %.2f (%s)\n", match.BestAwayWin.Value, match.BestAwayWin.SiteName)
			}
		}

		if len(bestOdds) > 5 {
			fmt.Printf("\n... and %d more matches\n", len(bestOdds)-5)
		}
	}

	fmt.Println("\n✨ Test completed!")
	
	if successfulSites == 0 {
		fmt.Println("\n⚠️  No sites were successfully scraped.")
		fmt.Println("This might be due to:")
		fmt.Println("- Network connectivity issues")
		fmt.Println("- Betting sites being down")
		fmt.Println("- Changes in website structure")
		fmt.Println("- Chrome/Chromium not being installed")
		log.Fatal("All scraping attempts failed")
	}
}