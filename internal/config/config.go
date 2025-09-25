package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                string
	ScrapeInterval      time.Duration
	MaxConcurrentScrapers int
	RequestTimeout      time.Duration
	ChromeHeadless      bool
	ChromeDisableGPU    bool
	RateLimitRequests   int
	RateLimitWindow     time.Duration
	LogLevel            string
}

func New() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		ScrapeInterval:      getDurationEnv("SCRAPE_INTERVAL", 300) * time.Second,
		MaxConcurrentScrapers: getIntEnv("MAX_CONCURRENT_SCRAPERS", 5),
		RequestTimeout:      getDurationEnv("REQUEST_TIMEOUT", 30) * time.Second,
		ChromeHeadless:      getBoolEnv("CHROME_HEADLESS", true),
		ChromeDisableGPU:    getBoolEnv("CHROME_DISABLE_GPU", true),
		RateLimitRequests:   getIntEnv("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:     getDurationEnv("RATE_LIMIT_WINDOW", 60) * time.Second,
		LogLevel:            getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return time.Duration(intValue)
		}
	}
	return time.Duration(defaultValue)
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}