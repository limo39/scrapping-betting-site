# Kenya Betting Odds Scraper 🏆

A powerful Go-based web scraper that collects and compares betting odds from popular Kenyan betting sites to help you find the best odds for football matches. Built with modern web technologies and designed for reliability and performance.

## ✨ Key Features

- 🔄 **Automated Scraping** - Scrapes odds every 5 minutes with smart scheduling
- 🏆 **Best Odds Comparison** - Intelligent comparison across all betting sites
- 🌐 **Modern Web Interface** - Clean, responsive UI with real-time updates
- 📊 **REST API** - Comprehensive JSON API for programmatic access
- ⚡ **High Performance** - Concurrent scraping with configurable limits
- 🛡️ **Rate Limiting** - Respectful scraping with built-in delays and error handling
- 📱 **Mobile Optimized** - Fully responsive design for all devices
- 🚀 **Demo Mode** - Fast testing without Chrome dependencies
- 🐳 **Docker Ready** - Complete containerization support

## 🎯 Supported Betting Sites

| Site | URL | Status |
|------|-----|--------|
| **Betika** | https://www.betika.com | ✅ Active |
| **SportPesa** | https://www.sportpesa.com | ✅ Active |
| **Betway** | https://www.betway.co.ke | ✅ Active |
| **Odibets** | https://www.odibets.com | ✅ Active |

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Chrome/Chromium** - For headless scraping (auto-installed in Docker)
- **Git** - For cloning the repository

### ⚡ Super Quick Demo (30 seconds)

```bash
# Clone and test immediately
git clone <repository-url>
cd scrapping-betting-site
make test-demo    # Test without Chrome dependencies
make demo         # Start demo server
# Visit http://localhost:8081 🎉
```

### 📦 Full Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd scrapping-betting-site
```

2. **One-command setup:**
```bash
make setup        # Installs dependencies and creates config
```

3. **Choose your mode:**

**Demo Mode** (Recommended for testing):
```bash
make demo         # Fast, no Chrome needed
# Visit http://localhost:8081
```

**Production Mode** (Real scraping):
```bash
make run          # Full scraping with Chrome
# Visit http://localhost:8080
```

## ⚙️ Configuration

The scraper supports multiple configuration modes:

### Environment Files

| File | Purpose | Usage |
|------|---------|-------|
| `.env.example` | Template | Copy to `.env` for custom config |
| `.env.demo` | Demo mode | Fast testing without Chrome |
| `.env` | Production | Your custom configuration |

### Key Settings

```env
# Server Configuration
PORT=8080                    # Server port
SCRAPE_INTERVAL=300         # Scraping interval in seconds (5 minutes)
MAX_CONCURRENT_SCRAPERS=5   # Max concurrent scrapers

# Performance
REQUEST_TIMEOUT=30          # Request timeout in seconds
RATE_LIMIT_REQUESTS=100     # Requests per window
RATE_LIMIT_WINDOW=60        # Rate limit window in seconds

# Chrome Settings
CHROME_HEADLESS=true        # Run Chrome in headless mode
CHROME_DISABLE_GPU=true     # Disable GPU acceleration

# Modes
LOG_LEVEL=info             # info, debug, demo
```

### Quick Configuration

```bash
# Demo mode (fast testing)
make demo

# Custom port
PORT=9000 make run

# Debug mode
LOG_LEVEL=debug make run
```

## API Endpoints

### Get Best Odds
```http
GET /api/v1/odds/best
```
Returns the best odds found across all betting sites.

### Trigger Manual Scrape
```http
POST /api/v1/scrape/trigger
```
Manually triggers a scraping operation.

### Get Scrape Results
```http
GET /api/v1/scrape/results
```
Returns the history of scraping operations.

### Health Check
```http
GET /api/v1/health
```
Returns service health status.

## Usage Examples

### Using the Web Interface

1. Visit http://localhost:8080
2. Click "Refresh Odds" to trigger scraping
3. View the best odds comparison in the cards
4. Odds are color-coded with the best odds highlighted

### Using the API

```bash
# Get best odds
curl http://localhost:8080/api/v1/odds/best

# Trigger scraping
curl -X POST http://localhost:8080/api/v1/scrape/trigger

# Check health
curl http://localhost:8080/api/v1/health
```

## Development

### Project Structure

```
scrapping-betting-site/
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # REST API handlers
│   ├── config/            # Configuration management
│   ├── models/            # Data models
│   ├── scraper/           # Scraping logic
│   └── scheduler/         # Cron scheduling
├── web/
│   └── templates/         # HTML templates
└── README.md
```

### Adding New Betting Sites

1. Create a new scraper in `internal/scraper/`
2. Implement the `Scraper` interface
3. Register the scraper in `manager.go`

Example:
```go
type NewSiteScraper struct {
    siteInfo models.BettingSite
}

func (n *NewSiteScraper) GetSiteInfo() models.BettingSite {
    return n.siteInfo
}

func (n *NewSiteScraper) ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error) {
    // Implement scraping logic
}
```

## Important Notes

### Legal and Ethical Considerations

- ⚖️ **Respect Terms of Service** - Always check and comply with each site's ToS
- 🤝 **Rate Limiting** - Built-in delays to avoid overwhelming servers
- 📋 **Personal Use** - This tool is for personal odds comparison only
- 🚫 **No Automated Betting** - Do not use for automated betting systems

### Technical Considerations

- 🌐 **Dynamic Content** - Uses Chrome headless browser for JavaScript-heavy sites
- 🔄 **Selector Updates** - Website changes may require selector updates
- 📊 **Data Accuracy** - Always verify odds on the actual betting site
- 🛡️ **Error Handling** - Graceful handling of site unavailability

## Troubleshooting

### Common Issues

1. **Chrome not found**: Install Chrome/Chromium browser
2. **Scraping fails**: Check if betting sites have changed their structure
3. **No odds data**: Ensure sites are accessible and selectors are correct
4. **High memory usage**: Reduce `MAX_CONCURRENT_SCRAPERS` in config

### Debugging

Enable debug logging:
```env
LOG_LEVEL=debug
```

Check scrape results:
```bash
curl http://localhost:8080/api/v1/scrape/results
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

This tool is for educational and personal use only. Always verify odds on the official betting sites before placing any bets. The developers are not responsible for any losses incurred from using this tool.