package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mileusna/useragent"
)

type ParsedAnalytics struct {
	DeviceType  string
	OS          string
	Browser     string
	Referrer    string
	IPAddress   string
	Country     string
	CountryCode string
	City        string
}

// ParseUserAgent extracts device type, OS, and browser from a user-agent string
func ParseUserAgent(uaString string) (deviceType, osName, browser string) {
	if uaString == "" {
		return "Unknown", "Unknown", "Unknown"
	}

	ua := useragent.Parse(uaString)

	// Device Type
	if ua.Mobile {
		deviceType = "Mobile"
	} else if ua.Tablet {
		deviceType = "Tablet"
	} else if ua.Bot {
		deviceType = "Bot"
	} else {
		deviceType = "Desktop"
	}

	// OS
	osName = ua.OS
	if osName == "" {
		osName = "Unknown OS"
	}

	// Browser
	browser = ua.Name
	if browser == "" {
		browser = "Unknown Browser"
	}

	return deviceType, osName, browser
}

// CleanReferrer parses and cleans the referrer domain or categorizes it
func CleanReferrer(rawReferrer string) string {
	if rawReferrer == "" {
		return "Direct"
	}

	parsed, err := url.Parse(rawReferrer)
	if err != nil || parsed.Host == "" {
		return "Direct"
	}

	host := strings.ToLower(parsed.Host)
	host = strings.TrimPrefix(host, "www.")

	// Categorize popular referrers
	switch {
	case strings.Contains(host, "google"):
		return "Google"
	case strings.Contains(host, "bing"):
		return "Bing"
	case strings.Contains(host, "yahoo"):
		return "Yahoo"
	case strings.Contains(host, "duckduckgo"):
		return "DuckDuckGo"
	case strings.Contains(host, "twitter") || strings.Contains(host, "t.co") || strings.Contains(host, "x.com"):
		return "Twitter / X"
	case strings.Contains(host, "facebook") || strings.Contains(host, "fb.com"):
		return "Facebook"
	case strings.Contains(host, "instagram"):
		return "Instagram"
	case strings.Contains(host, "linkedin"):
		return "LinkedIn"
	case strings.Contains(host, "youtube"):
		return "YouTube"
	case strings.Contains(host, "tiktok"):
		return "TikTok"
	case strings.Contains(host, "reddit"):
		return "Reddit"
	case strings.Contains(host, "whatsapp"):
		return "WhatsApp"
	case strings.Contains(host, "telegram"):
		return "Telegram"
	default:
		return host
	}
}

// ExtractClientIP retrieves client IP from request headers or remote address
func ExtractClientIP(r *http.Request) string {
	// Check standard proxy headers
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"} {
		val := r.Header.Get(header)
		if val != "" {
			parts := strings.Split(val, ",")
			ip := strings.TrimSpace(parts[0])
			if ip != "" && !isPrivateIP(ip) {
				return ip
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}

type ipApiResponse struct {
	Status      string `json:"status"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
}

// LookupGeoLocation resolves IP to Country and City (with fast timeout and local fallback)
func LookupGeoLocation(ip string) (country, countryCode, city string) {
	if ip == "" || ip == "127.0.0.1" || ip == "::1" || isPrivateIP(ip) {
		return "Localhost", "LOCAL", "Local Development"
	}

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,city", ip))
	if err != nil {
		return "Unknown", "XX", "Unknown"
	}
	defer resp.Body.Close()

	var geo ipApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil || geo.Status != "success" {
		return "Unknown", "XX", "Unknown"
	}

	return geo.Country, geo.CountryCode, geo.City
}
