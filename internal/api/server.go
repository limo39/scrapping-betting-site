package api

import (
	"context"
	"net/http"
	"time"

	"betting-odds-scraper/internal/scraper"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	manager *scraper.Manager
}

func NewServer(manager *scraper.Manager) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		router:  router,
		manager: manager,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := s.router.Group("/api/v1")
	{
		api.GET("/health", s.healthCheck)
		api.GET("/odds/best", s.getBestOdds)
		api.GET("/odds/stats", s.getOddsStats)
		api.GET("/scrape/results", s.getScrapeResults)
		api.POST("/scrape/trigger", s.triggerScrape)
		api.GET("/sites", s.getSites)
		api.GET("/sites/status", s.getSitesStatus)
	}

	// Serve static files for simple web interface
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("web/templates/*")
	s.router.GET("/", s.indexPage)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "betting-odds-scraper",
	})
}

func (s *Server) getBestOdds(c *gin.Context) {
	bestOdds := s.manager.GetBestOdds()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bestOdds,
		"count":   len(bestOdds),
	})
}

func (s *Server) getScrapeResults(c *gin.Context) {
	results := s.manager.GetScrapeResults()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}

func (s *Server) triggerScrape(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	results := s.manager.ScrapeAll(ctx)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Scraping completed",
		"results": results,
	})
}

func (s *Server) getSites(c *gin.Context) {
	// This would return information about supported betting sites
	sites := []map[string]interface{}{
		{"id": "betika", "name": "Betika", "active": true},
		{"id": "sportpesa", "name": "SportPesa", "active": true},
		{"id": "betway", "name": "Betway", "active": true},
		{"id": "odibets", "name": "Odibets", "active": true},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sites,
	})
}

func (s *Server) getOddsStats(c *gin.Context) {
	bestOdds := s.manager.GetBestOdds()
	
	totalMatches := len(bestOdds)
	totalSites := 4 // Betika, SportPesa, Betway, Odibets
	
	// Calculate average odds and other stats
	var totalHomeOdds, totalDrawOdds, totalAwayOdds float64
	var homeCount, drawCount, awayCount int
	
	for _, match := range bestOdds {
		if match.BestHomeWin != nil {
			totalHomeOdds += match.BestHomeWin.Value
			homeCount++
		}
		if match.BestDraw != nil {
			totalDrawOdds += match.BestDraw.Value
			drawCount++
		}
		if match.BestAwayWin != nil {
			totalAwayOdds += match.BestAwayWin.Value
			awayCount++
		}
	}
	
	stats := gin.H{
		"total_matches": totalMatches,
		"total_sites": totalSites,
		"average_home_odds": func() float64 {
			if homeCount > 0 { return totalHomeOdds / float64(homeCount) }
			return 0
		}(),
		"average_draw_odds": func() float64 {
			if drawCount > 0 { return totalDrawOdds / float64(drawCount) }
			return 0
		}(),
		"average_away_odds": func() float64 {
			if awayCount > 0 { return totalAwayOdds / float64(awayCount) }
			return 0
		}(),
		"last_updated": time.Now(),
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": stats,
	})
}

func (s *Server) getSitesStatus(c *gin.Context) {
	results := s.manager.GetScrapeResults()
	
	sites := []gin.H{
		{"id": "betika", "name": "Betika", "url": "https://www.betika.com"},
		{"id": "sportpesa", "name": "SportPesa", "url": "https://www.sportpesa.com"},
		{"id": "betway", "name": "Betway", "url": "https://www.betway.co.ke"},
		{"id": "odibets", "name": "Odibets", "url": "https://www.odibets.com"},
	}
	
	// Add status information from latest scrape results
	for i, site := range sites {
		siteID := site["id"].(string)
		if siteResults, exists := results[siteID]; exists && len(siteResults) > 0 {
			latest := siteResults[len(siteResults)-1]
			sites[i]["status"] = map[string]interface{}{
				"active": latest.Success,
				"last_scrape": latest.ScrapedAt,
				"match_count": latest.MatchCount,
				"odds_count": latest.OddsCount,
				"error": latest.Error,
			}
		} else {
			sites[i]["status"] = map[string]interface{}{
				"active": false,
				"last_scrape": nil,
				"match_count": 0,
				"odds_count": 0,
				"error": "No data available",
			}
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": sites,
	})
}

func (s *Server) indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Kenya Betting Odds Scraper",
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}