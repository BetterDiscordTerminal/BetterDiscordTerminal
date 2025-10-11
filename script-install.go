package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var TMP = "/tmp/betterdiscord.asar"

// RunInstallation performs the BetterDiscord installation silently.
func RunInstallation(discordVersion DiscordVersion) error {
	if err := KillDiscord(discordVersion); err != nil {
		return fmt.Errorf("failed to kill Discord: %w", err)
	}

	if err := CreateBDDirectory(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	asarURL, err := GetAsarURL()
	if err != nil {
		return fmt.Errorf("failed to get BetterDiscord.asar URL: %w", err)
	}

	if err := DownloadASAR(asarURL); err != nil {
		return fmt.Errorf("failed to download BetterDiscord.asar: %w", err)
	}

	bdAsarPath, err := MoveBetterDiscordAsar()
	if err != nil {
		return fmt.Errorf("failed to move BetterDiscord.asar: %w", err)
	}

	if err := InjectShim(bdAsarPath, discordVersion); err != nil {
		return fmt.Errorf("failed to inject shim: %w", err)
	}

	LaunchDiscord(discordVersion)
	return nil
}

// CreateBDDirectory creates the BetterDiscord directory structure.
func CreateBDDirectory() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	bdPath := filepath.Join(home, "Library", "Application Support", "BetterDiscord")
	subdirs := []string{"", "data", "plugins", "themes"}

	for _, dir := range subdirs {
		dirPath := filepath.Join(bdPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return nil
}

// Gets the .asar download url from github api.
func GetAsarURL() (string, error) {
	url := "https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest"
	client := http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "betterdiscord-term")
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release: %w", err)
	}
	defer res.Body.Close()

	release := Release{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	respByte := buf.Bytes()

	if err := json.Unmarshal(respByte, &release); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(release.Assets) == 0 || release.Assets[0].Name != "betterdiscord.asar" {
		return "", fmt.Errorf("BetterDiscord asset not found in latest release")
	}

	return release.Assets[0].BrowserDownloadURL, nil
}

// Downloads the asar file from the given URL and stores it in /tmp directory.
func DownloadASAR(asarPath string) error {
	client := http.Client{Timeout: time.Second * 30}

	asarReq, err := http.NewRequest(http.MethodGet, asarPath, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	asarReq.Header.Set("User-Agent", "betterdiscord-term")
	asarRes, err := client.Do(asarReq)
	if err != nil {
		return fmt.Errorf("failed to download asar: %w", err)
	}
	defer asarRes.Body.Close()

	if asarRes.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download asar, status: %d", asarRes.StatusCode)
	}

	outFile, err := os.Create(TMP)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, asarRes.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// Move BetterDiscord.asar from /tmp to the BetterDiscord data directory.
func MoveBetterDiscordAsar() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	destPath := filepath.Join(home, "Library", "Application Support", "BetterDiscord", "data", "betterdiscord.asar")
	if err := os.Rename(TMP, destPath); err != nil {
		return "", fmt.Errorf("failed to move asar file: %w", err)
	}

	return destPath, nil
}

// Change index.js file to load BetterDiscord asar.
func InjectShim(asarPath string, discordVersion DiscordVersion) error {
	shimContent := fmt.Sprintf(`require('%s');
module.exports = require("./core.asar");`, asarPath)

	corePath, err := DiscordApplicationSupportPath(discordVersion)
	if err != nil {
		return fmt.Errorf("failed to get Discord path: %w", err)
	}

	shimPath := filepath.Join(corePath, "index.js")
	if err := os.WriteFile(shimPath, []byte(shimContent), 0644); err != nil {
		return fmt.Errorf("failed to write shim: %w", err)
	}

	return nil
}
