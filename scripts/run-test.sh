#!/bin/bash

# Kenya Betting Odds Scraper Test Script
echo "🏆 Kenya Betting Odds Scraper Test"
echo "=================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or higher."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Check if Chrome/Chromium is installed
if command -v google-chrome &> /dev/null; then
    echo "✅ Google Chrome detected"
elif command -v chromium-browser &> /dev/null; then
    echo "✅ Chromium browser detected"
elif command -v chromium &> /dev/null; then
    echo "✅ Chromium detected"
else
    echo "⚠️  Chrome/Chromium not found. Installing..."
    
    # Try to install Chrome/Chromium based on OS
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y chromium-browser
        elif command -v yum &> /dev/null; then
            sudo yum install -y chromium
        elif command -v pacman &> /dev/null; then
            sudo pacman -S chromium
        else
            echo "❌ Please install Chrome or Chromium manually"
            exit 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew install --cask google-chrome
        else
            echo "❌ Please install Chrome manually or install Homebrew"
            exit 1
        fi
    else
        echo "❌ Unsupported OS. Please install Chrome/Chromium manually"
        exit 1
    fi
fi

# Navigate to project directory
cd "$(dirname "$0")/.."

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
fi

# Run the test
echo "🧪 Running scraper test..."
go run scripts/test-scraper.go

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 Test completed successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Run 'make run' or 'go run main.go' to start the server"
    echo "2. Open http://localhost:8080 in your browser"
    echo "3. Click 'Refresh Odds' to scrape live data"
    echo ""
    echo "API endpoints:"
    echo "- GET  /api/v1/odds/best     - Get best odds"
    echo "- POST /api/v1/scrape/trigger - Trigger scraping"
    echo "- GET  /api/v1/health        - Health check"
else
    echo "❌ Test failed. Check the error messages above."
    exit 1
fi