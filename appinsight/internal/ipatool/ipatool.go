package ipatool

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type SearchResult struct {
	Name      string                 `json:"name"`
	BundleID  string                 `json:"bundleId"`
	Version   string                 `json:"version"`
	Developer string                 `json:"developer,omitempty"`
	TrackID   int64                  `json:"trackId"`
	Price     interface{}            `json:"price"`
	Raw       map[string]interface{} `json:"raw"`
}

type SearchResponse struct {
	OK      bool           `json:"ok"`
	Command string         `json:"command"`
	Count   int            `json:"count"`
	Results []SearchResult `json:"results"`
}

type ipatoolSearchResponse struct {
	Level string `json:"level"`
	Count int    `json:"count"`
	Apps  []struct {
		ID        int64       `json:"id"`
		BundleID  string      `json:"bundleID"`
		Name      string      `json:"name"`
		Version   string      `json:"version"`
		Price     interface{} `json:"price"`
	} `json:"apps"`
}

func Search(keyword string, limit int) (*SearchResponse, error) {
	args := []string{"search", keyword, "--format", "json"}
	if limit > 0 {
		args = append(args, "--limit", fmt.Sprintf("%d", limit))
	}

	cmd := exec.Command("ipatool", args...)
	stdout, err := cmd.Output()
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("ipatool is not installed. Install it via: brew install majd/repo/ipatool")
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("ipatool search failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("ipatool search failed: %w", err)
	}

	var rawResp ipatoolSearchResponse
	if err := json.Unmarshal(stdout, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to parse ipatool output as JSON: %w\nraw output: %s", err, string(stdout))
	}

	results := make([]SearchResult, 0, len(rawResp.Apps))
	for _, app := range rawResp.Apps {
		raw := map[string]interface{}{
			"id":        app.ID,
			"bundleID":  app.BundleID,
			"name":      app.Name,
			"version":   app.Version,
			"price":     app.Price,
		}
		results = append(results, SearchResult{
			Name:     app.Name,
			BundleID: app.BundleID,
			Version:  app.Version,
			TrackID:  app.ID,
			Price:    app.Price,
			Raw:      raw,
		})
	}

	return &SearchResponse{
		OK:      true,
		Command: "search",
		Count:   len(results),
		Results: results,
	}, nil
}

type FetchResponse struct {
	OK                bool   `json:"ok"`
	Command           string `json:"command"`
	BundleID          string `json:"bundleId"`
	IPAPath           string `json:"ipaPath"`
	FileSize          int64  `json:"fileSize"`
	SHA256            string `json:"sha256"`
	EncryptedReminder string `json:"encryptedReminder"`
}

type ipatoolDownloadResponse struct {
	Level  string `json:"level"`
	Path   string `json:"path"`
}

func Fetch(bundleID, outputDir string, purchase bool) (*FetchResponse, error) {
	args := []string{"download", "--bundle-identifier", bundleID}
	if purchase {
		args = append(args, "--purchase")
	}
	if outputDir != "" {
		args = append(args, "--output", outputDir)
	}
	args = append(args, "--format", "json")

	cmd := exec.Command("ipatool", args...)
	stdout, err := cmd.Output()
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("ipatool is not installed. Install it via: brew install majd/repo/ipatool")
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "not authenticated") || strings.Contains(stderr, "login") || strings.Contains(stderr, "auth") {
				return nil, fmt.Errorf("ipatool is not logged in. Please run: ipatool auth login")
			}
			return nil, fmt.Errorf("ipatool download failed: %s", stderr)
		}
		return nil, fmt.Errorf("ipatool download failed: %w", err)
	}

	var dlOutput ipatoolDownloadResponse
	if err := json.Unmarshal(stdout, &dlOutput); err != nil {
		return nil, fmt.Errorf("failed to parse ipatool download output: %w\nraw output: %s", err, string(stdout))
	}

	ipaPath := dlOutput.Path
	fileSize := int64(0)
	sha256 := ""

	if info, err := exec.Command("stat", "-f%z", ipaPath).Output(); err == nil {
		fmt.Sscanf(strings.TrimSpace(string(info)), "%d", &fileSize)
	}

	if hash, err := exec.Command("shasum", "-a", "256", ipaPath).Output(); err == nil {
		parts := strings.Split(string(hash), " ")
		if len(parts) > 0 {
			sha256 = parts[0]
		}
	}

	return &FetchResponse{
		OK:                true,
		Command:           "fetch-ios",
		BundleID:          bundleID,
		IPAPath:           ipaPath,
		FileSize:          fileSize,
		SHA256:            sha256,
		EncryptedReminder: "App Store IPA is usually encrypted. Only visible structure analysis is available.",
	}, nil
}

func isNotFound(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return strings.Contains(string(exitErr.Stderr), "not found") || strings.Contains(exitErr.Error(), "executable file not found")
	}
	return strings.Contains(err.Error(), "executable file not found") || strings.Contains(err.Error(), "not found")
}
