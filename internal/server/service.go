package server

import (
	"encoding/json"
	"io"
	"license/internal/logger"
	"license/internal/useragent"
	"net/http"
	"strings"
	"time"
)

// VersionResponse defines the version response structure
type VersionResponse struct {
	Version       string `json:"version"`
	Build         string `json:"build"`
	OsArch        string `json:"osArch"`
	NeedUpdate    bool   `json:"needUpdate"`
	LatestVersion string `json:"latestVersion,omitempty"`
}

// GitHubRelease GitHub API release response structure
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// Cache for GitHub latest version information
var (
	cachedLatestVersion string
	lastFetchTime       time.Time
	cacheExpiration     = 5 * time.Minute
)

// getLatestVersionFromGitHub fetches the latest version from GitHub(removing the "v" prefix)
func getLatestVersionFromGitHub() string {
	// If the cache has not expired, use the cached value
	if !lastFetchTime.IsZero() && time.Since(lastFetchTime) < cacheExpiration && cachedLatestVersion != "" {
		return cachedLatestVersion
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/nannanStrawberry314/license/releases/latest", nil)
	if err != nil {
		logger.Error("Failed to create GitHub API request", err)
		return ""
	}

	// Use random User-Agent
	req.Header.Set("User-Agent", useragent.GetRandom())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to request GitHub API", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("GitHub API returned non-200 status code", nil)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read GitHub API response", err)
		return ""
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Error("Failed to parse GitHub API response", err)
		return ""
	}

	// Remove "v" prefix from version number
	version := release.TagName
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	// Update cache
	cachedLatestVersion = version
	lastFetchTime = time.Now()

	return version
}

// compareVersions compares version numbers
func compareVersions(current, latest string) bool {
	// Split both versions by dot
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// Compare each part of the version numbers
	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		if currentParts[i] < latestParts[i] {
			return true // Update needed
		} else if currentParts[i] > latestParts[i] {
			return false // No update needed
		}
	}

	// If all previous parts are equal, but latest has more parts, update is needed
	return len(latestParts) > len(currentParts)
}
