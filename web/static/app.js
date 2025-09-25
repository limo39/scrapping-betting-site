// Additional JavaScript functionality for Kenya Betting Odds Scraper

// Service Worker for offline support (future enhancement)
if ('serviceWorker' in navigator) {
    window.addEventListener('load', function() {
        // navigator.serviceWorker.register('/sw.js');
    });
}

// Keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl/Cmd + R for refresh
    if ((e.ctrlKey || e.metaKey) && e.key === 'r') {
        e.preventDefault();
        triggerScrape();
    }
    
    // Escape to clear search
    if (e.key === 'Escape') {
        const searchInput = document.getElementById('search-input');
        if (searchInput.value) {
            searchInput.value = '';
            filterMatches();
        }
    }
    
    // Ctrl/Cmd + F to focus search
    if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
        e.preventDefault();
        document.getElementById('search-input').focus();
    }
});

// Enhanced error handling
window.addEventListener('error', function(e) {
    console.error('JavaScript error:', e.error);
    showToast('An unexpected error occurred', 'error');
});

// Network status monitoring
window.addEventListener('online', function() {
    showToast('Connection restored', 'success');
    if (autoRefreshEnabled) {
        loadOdds();
    }
});

window.addEventListener('offline', function() {
    showToast('Connection lost - working offline', 'warning');
});

// Performance monitoring
function measurePerformance(name, fn) {
    const start = performance.now();
    const result = fn();
    const end = performance.now();
    console.log(`${name} took ${end - start} milliseconds`);
    return result;
}

// Local storage for user preferences
function savePreference(key, value) {
    try {
        localStorage.setItem(`betting-scraper-${key}`, JSON.stringify(value));
    } catch (e) {
        console.warn('Could not save preference:', e);
    }
}

function getPreference(key, defaultValue) {
    try {
        const stored = localStorage.getItem(`betting-scraper-${key}`);
        return stored ? JSON.parse(stored) : defaultValue;
    } catch (e) {
        console.warn('Could not load preference:', e);
        return defaultValue;
    }
}

// Initialize user preferences
document.addEventListener('DOMContentLoaded', function() {
    // Restore auto-refresh preference
    const savedAutoRefresh = getPreference('autoRefresh', true);
    if (!savedAutoRefresh && autoRefreshEnabled) {
        toggleAutoRefresh();
    }
    
    // Restore last selected league filter
    const savedLeague = getPreference('selectedLeague', 'all');
    if (savedLeague !== 'all') {
        setTimeout(() => {
            const tab = document.querySelector(`[onclick="filterByLeague('${savedLeague}')"]`);
            if (tab) tab.click();
        }, 1000);
    }
});

// Save preferences when changed
function saveAutoRefreshPreference() {
    savePreference('autoRefresh', autoRefreshEnabled);
}

function saveLeaguePreference(league) {
    savePreference('selectedLeague', league);
}

// Enhanced analytics (privacy-friendly)
function trackEvent(category, action, label) {
    // This would integrate with privacy-friendly analytics
    console.log('Event:', category, action, label);
}

// Accessibility improvements
function announceToScreenReader(message) {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.className = 'sr-only';
    announcement.textContent = message;
    
    document.body.appendChild(announcement);
    
    setTimeout(() => {
        document.body.removeChild(announcement);
    }, 1000);
}

// Export functionality
function exportOddsData(format = 'json') {
    if (!oddsData.length) {
        showToast('No data to export', 'warning');
        return;
    }
    
    let content, filename, mimeType;
    
    switch (format) {
        case 'csv':
            content = convertToCSV(oddsData);
            filename = 'betting-odds.csv';
            mimeType = 'text/csv';
            break;
        case 'json':
        default:
            content = JSON.stringify(oddsData, null, 2);
            filename = 'betting-odds.json';
            mimeType = 'application/json';
            break;
    }
    
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
    
    showToast(`Data exported as ${format.toUpperCase()}`, 'success');
}

function convertToCSV(data) {
    const headers = ['Home Team', 'Away Team', 'League', 'Home Odds', 'Home Site', 'Draw Odds', 'Draw Site', 'Away Odds', 'Away Site', 'Updated'];
    const rows = data.map(match => [
        match.match.home_team,
        match.match.away_team,
        match.match.league,
        match.best_home_win?.value || '',
        match.best_home_win?.site_name || '',
        match.best_draw?.value || '',
        match.best_draw?.site_name || '',
        match.best_away_win?.value || '',
        match.best_away_win?.site_name || '',
        new Date(match.updated_at).toLocaleString()
    ]);
    
    return [headers, ...rows].map(row => 
        row.map(field => `"${field}"`).join(',')
    ).join('\n');
}

// Comparison tools
function compareOdds(match) {
    if (!match.all_odds || match.all_odds.length < 2) {
        showToast('Not enough data for comparison', 'warning');
        return;
    }
    
    // Create comparison modal or popup
    const comparison = match.all_odds.map(odd => ({
        site: odd.site_name,
        home: odd.home_win,
        draw: odd.draw,
        away: odd.away_win
    }));
    
    console.log('Odds comparison:', comparison);
    // This would show a detailed comparison modal
}

// Notification system
function requestNotificationPermission() {
    if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission().then(permission => {
            if (permission === 'granted') {
                showToast('Notifications enabled', 'success');
            }
        });
    }
}

function showNotification(title, body, icon = '/favicon.ico') {
    if ('Notification' in window && Notification.permission === 'granted') {
        new Notification(title, { body, icon });
    }
}

// Initialize additional features
document.addEventListener('DOMContentLoaded', function() {
    // Add keyboard shortcut hints
    const shortcuts = document.createElement('div');
    shortcuts.className = 'd-none d-md-block position-fixed bottom-0 start-0 p-2 text-muted small';
    shortcuts.innerHTML = `
        <div>Shortcuts: Ctrl+R (Refresh) | Ctrl+F (Search) | Esc (Clear)</div>
    `;
    document.body.appendChild(shortcuts);
    
    // Request notification permission after user interaction
    setTimeout(requestNotificationPermission, 5000);
});