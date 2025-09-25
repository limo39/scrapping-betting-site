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

## 📡 API Reference

### Core Endpoints

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/api/v1/odds/best` | Get best odds comparison | JSON with best odds across all sites |
| `POST` | `/api/v1/scrape/trigger` | Trigger manual scrape | Scraping results and status |
| `GET` | `/api/v1/scrape/results` | Get scrape history | Historical scraping data |
| `GET` | `/api/v1/health` | Health check | Service status and uptime |
| `GET` | `/api/v1/sites` | List supported sites | Available betting sites |

### Example Responses

**Best Odds:**
```json
{
  "success": true,
  "data": [
    {
      "match": {
        "home_team": "Arsenal",
        "away_team": "Chelsea",
        "league": "Premier League"
      },
      "best_home_win": {
        "value": 2.15,
        "site_name": "Betway"
      },
      "best_away_win": {
        "value": 3.25,
        "site_name": "SportPesa"
      }
    }
  ],
  "count": 24
}
```

**Health Check:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "betting-odds-scraper"
}
```

## 💡 Usage Examples

### 🌐 Web Interface

1. **Start the server:**
   ```bash
   make demo  # or make run
   ```

2. **Open your browser:**
   - Demo mode: http://localhost:8081
   - Production: http://localhost:8080

3. **Use the interface:**
   - Click "Refresh Odds" to trigger scraping
   - View best odds comparison in responsive cards
   - Auto-refresh every 5 minutes
   - Mobile-friendly design

### 🔧 API Usage

```bash
# Get best odds with formatting
curl -s http://localhost:8081/api/v1/odds/best | jq .

# Trigger manual scrape
curl -X POST http://localhost:8081/api/v1/scrape/trigger

# Check service health
curl http://localhost:8081/api/v1/health

# Get scraping history
curl http://localhost:8081/api/v1/scrape/results

# List supported sites
curl http://localhost:8081/api/v1/sites
```

### 🛠️ Command Line Tools

```bash
# Quick test all scrapers
make test-demo

# Test with real Chrome scraping
make test-simple

# Check service health
make health

# Get current best odds
make odds

# Manual scrape trigger
make scrape
```

## 🏗️ Development

### Project Architecture

```
scrapping-betting-site/
├── main.go                 # Application entry point
├── cmd/                   # Command-line tools
│   └── test-simple/       # Simple testing utility
├── internal/              # Private application code
│   ├── api/              # REST API handlers & server
│   ├── config/           # Configuration management
│   ├── models/           # Data structures & types
│   ├── scraper/          # Scraping engines
│   │   ├── manager.go    # Scraper orchestration
│   │   ├── demo.go       # Demo mode scraper
│   │   ├── betika.go     # Betika scraper
│   │   ├── sportpesa.go  # SportPesa scraper
│   │   ├── betway.go     # Betway scraper
│   │   └── odibets.go    # Odibets scraper
│   └── scheduler/        # Cron job scheduling
├── web/
│   └── templates/        # HTML templates
├── scripts/              # Utility scripts
├── Dockerfile           # Container configuration
├── docker-compose.yml   # Multi-container setup
└── Makefile            # Build automation
```

### 🔧 Available Make Commands

```bash
# Development
make setup          # First-time setup
make run            # Start production server
make demo           # Start demo server
make dev            # Start with live reload (requires air)

# Testing
make test           # Run unit tests
make test-demo      # Quick demo test
make test-simple    # Test without Chrome noise
make test-coverage  # Test with coverage report

# Building
make build          # Build binary
make clean          # Clean build artifacts
make install        # Install to system

# Docker
make docker-build   # Build Docker image
make docker-run     # Run in container

# Utilities
make fmt            # Format code
make lint           # Lint code (requires golangci-lint)
make security       # Security scan (requires gosec)
make health         # Check service health
make odds           # Get current best odds
make scrape         # Trigger manual scrape
```

### 🆕 Adding New Betting Sites

1. **Create scraper file:**
   ```bash
   touch internal/scraper/newsite.go
   ```

2. **Implement the interface:**
   ```go
   type NewSiteScraper struct {
       siteInfo models.BettingSite
   }

   func NewNewSiteScraper() *NewSiteScraper {
       return &NewSiteScraper{
           siteInfo: models.BettingSite{
               ID:     "newsite",
               Name:   "New Site",
               URL:    "https://www.newsite.com",
               Active: true,
           },
       }
   }

   func (n *NewSiteScraper) GetSiteInfo() models.BettingSite {
       return n.siteInfo
   }

   func (n *NewSiteScraper) ScrapeOdds(ctx context.Context) ([]models.Match, []models.Odds, error) {
       // Implement scraping logic with Chrome or HTTP client
       // Return matches and odds data
   }
   ```

3. **Register in manager:**
   ```go
   // In internal/scraper/manager.go
   manager.RegisterScraper(NewNewSiteScraper())
   ```

### 🧪 Testing Your Changes

```bash
# Test new scraper
make test-demo

# Test with real scraping
make test-simple

# Full integration test
make run
curl http://localhost:8080/api/v1/odds/best
```

## 🚨 Important Considerations

### ⚖️ Legal & Ethical Guidelines

| ⚠️ **IMPORTANT** | **Guidelines** |
|------------------|----------------|
| **Terms of Service** | Always review and comply with each betting site's ToS |
| **Rate Limiting** | Built-in delays prevent server overload (respectful scraping) |
| **Personal Use Only** | This tool is for personal odds comparison, not commercial use |
| **No Auto-Betting** | Never use for automated betting or gambling systems |
| **Data Verification** | Always verify odds on official sites before placing bets |

### 🔧 Technical Considerations

| Aspect | Details |
|--------|---------|
| **Dynamic Content** | Uses Chrome headless for JavaScript-heavy sites |
| **Selector Maintenance** | Website changes require CSS selector updates |
| **Error Handling** | Graceful fallbacks for site unavailability |
| **Performance** | Concurrent scraping with configurable limits |
| **Reliability** | Built-in retry mechanisms and timeout handling |

### 🛡️ Security Features

- **No credentials stored** - Read-only public data access
- **Rate limiting** - Prevents overwhelming target servers  
- **Error isolation** - Failed scrapers don't affect others
- **Timeout protection** - Prevents hanging requests
- **Input validation** - Sanitized data processing

## 🔍 Troubleshooting

### Common Issues & Solutions

| Issue | Symptoms | Solution |
|-------|----------|----------|
| **Chrome not found** | `chrome: not found` error | Install Chrome/Chromium or use demo mode |
| **Port already in use** | `bind: address already in use` | Change PORT in `.env` or kill existing process |
| **No odds data** | Empty results, 0 matches | Use demo mode or update CSS selectors |
| **High memory usage** | System slowdown | Reduce `MAX_CONCURRENT_SCRAPERS` |
| **Timeout errors** | Context deadline exceeded | Increase `REQUEST_TIMEOUT` |
| **Permission denied** | Docker/Chrome issues | Add `--no-sandbox` flag or run as root |

### 🐛 Debugging Steps

1. **Enable debug logging:**
   ```bash
   LOG_LEVEL=debug make run
   ```

2. **Test individual components:**
   ```bash
   make test-demo      # Test without Chrome
   make health         # Check service status
   make scrape         # Manual scrape test
   ```

3. **Check scrape results:**
   ```bash
   curl http://localhost:8081/api/v1/scrape/results | jq .
   ```

4. **Monitor logs:**
   ```bash
   # In demo mode
   make demo

   # Check specific scraper
   LOG_LEVEL=debug make test-simple
   ```

### 🚀 Performance Optimization

```bash
# Fast demo mode (no Chrome)
make demo

# Reduce concurrent scrapers
MAX_CONCURRENT_SCRAPERS=2 make run

# Increase timeout for slow sites
REQUEST_TIMEOUT=60 make run

# Disable GPU for better Docker performance
CHROME_DISABLE_GPU=true make run
```

### 🆘 Getting Help

1. **Check the logs** - Most issues are logged with clear error messages
2. **Try demo mode** - Isolates Chrome/network issues
3. **Review FIXES.md** - Common issues and their solutions
4. **Test API directly** - Bypass web interface issues

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### 🔄 Development Workflow

1. **Fork & Clone:**
   ```bash
   git fork <repository-url>
   git clone <your-fork-url>
   cd scrapping-betting-site
   ```

2. **Setup Development Environment:**
   ```bash
   make setup
   make test-demo  # Verify everything works
   ```

3. **Create Feature Branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make Changes & Test:**
   ```bash
   make test-demo      # Quick test
   make test-coverage  # Full test suite
   make lint          # Code quality
   ```

5. **Submit Pull Request:**
   - Clear description of changes
   - Include tests for new features
   - Update documentation if needed

### 🎯 Contribution Ideas

- **New Betting Sites** - Add more Kenyan betting platforms
- **Enhanced Selectors** - Improve real website scraping
- **Data Storage** - Add database persistence
- **Mobile App** - React Native or Flutter frontend
- **Analytics** - Historical odds tracking and analysis
- **Notifications** - Alert system for odds changes

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

> **IMPORTANT**: This tool is for **educational and personal use only**.
> 
> - ✅ **Use for**: Learning web scraping, odds comparison, personal research
> - ❌ **Do not use for**: Commercial purposes, automated betting, ToS violations
> - 🔍 **Always verify**: Check odds on official betting sites before placing bets
> - 🚫 **No liability**: Developers are not responsible for any losses or issues
> 
> **Gambling can be addictive. Please gamble responsibly.**

---

## 🌟 Show Your Support

If this project helped you, please consider:
- ⭐ **Starring the repository**
- 🐛 **Reporting bugs or issues**
- 💡 **Suggesting new features**
- 🤝 **Contributing code or documentation**

**Happy scraping! 🚀**