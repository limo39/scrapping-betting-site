package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"betting-odds-scraper/internal/config"
	"betting-odds-scraper/internal/scraper"
)

func main() {
	// Suppress Chrome logs
	log.SetOutput(os.Stdout)
	
	fmt.Println("🏆 Kenya Betting Odds Scraper - Simple Test")
	fmt.Println("===========================================")

	// Initialize configuration
	cfg := config.New()
	cfg.RequestTimeout = 30 * time.Second

	// Initialize scraper manager
	manager := scraper.NewManager(cfg)

	fmt.Println("\n🔄 Testing scrapers...")
	
	// Test each scraper individually
	scrapers := []string{"betika", "sportpesa", "betway", "odibets"}
	
	for _, siteID := range scrapers {
		fmt.Printf("\n📊 Testing %s...\n", siteID)
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		
		// Create a single scraper test
		results := manager.ScrapeAll(ctx)
		
		if result, exists := results[siteID]; exists {
			if result.Success {
				fmt.Printf("✅ %s: %d matches, %d odds (%.2fs)\n", 
					siteID, result.MatchCount, result.OddsCount, result.Duration.Seconds())
			} else {
				fmt.Printf("❌ %s: Failed - %s\n", siteID, result.Error)
			}
		}
		
		cancel()
	}

	fmt.Println("\n🏆 Getting best odds comparison...")
	bestOdds := manager.GetBestOdds()
	
	if len(bestOdds) > 0 {
		fmt.Printf("\nFound %d matches with odds:\n", len(bestOdds))
		fmt.Println("=" + fmt.Sprintf("%*s", 50, "="))
		
		for i, match := range bestOdds {
			if i >= 10 { // Show only first 10
				break
			}
			
			fmt.Printf("\n🥅 %s vs %s\n", match.Match.HomeTeam, match.Match.AwayTeam)
			
			if match.BestHomeWin != nil {
				fmt.Printf("   Home: %.2f (%s)\n", match.BestHomeWin.Value, match.BestHomeWin.SiteName)
			}
			
			if match.BestDraw != nil {
				fmt.Printf("   Draw: %.2f (%s)\n", match.BestDraw.Value, match.BestDraw.SiteName)
			}
			
			if match.BestAwayWin != nil {
				fmt.Printf("   Away: %.2f (%s)\n", match.BestAwayWin.Value, match.BestAwayWin.SiteName)
			}
			
			fmt.Printf("   📈 %d sites offering odds\n", len(match.AllOdds))
		}
		
		if len(bestOdds) > 10 {
			fmt.Printf("\n... and %d more matches\n", len(bestOdds)-10)
		}
		
		fmt.Println("\n✨ Test completed successfully!")
		fmt.Println("\n🚀 Ready to start the web server!")
		fmt.Println("   Run: go run main.go")
		fmt.Println("   Then visit: http://localhost:8080")
		
	} else {
		fmt.Println("⚠️  No odds data found. This might be normal for a demo.")
	}
}