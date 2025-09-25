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
		api.GET("/scrape/results", s.getScrapeResults)
		api.POST("/scrape/trigger", s.triggerScrape)
		api.GET("/sites", s.getSites)
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

func (s *Server) indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Betting Odds Scraper",
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}