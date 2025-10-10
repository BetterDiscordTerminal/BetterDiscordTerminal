package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// killDiscordProcesses attempts to terminate all running instances of the specified Discord process.
func KillDiscord(discordVersion DiscordVersion) error {
	var processName string

	switch discordVersion.Index {
	case 0:
		processName = "Discord"
	case 1:
		processName = "Discord PTB"
	case 2:
		processName = "Discord Canary"
	default:
		return fmt.Errorf("unknown Discord version index: %d", discordVersion.Index)
	}

	// Check if process is running
	cmd := exec.Command("pgrep", "-x", processName)
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) == "" {
		return nil
	}

	killCmd := exec.Command("killall", processName)
	if err := killCmd.Run(); err != nil {
		return fmt.Errorf("failed to stop %s: %w", discordVersion.Name, err)
	}

	time.Sleep(2 * time.Second)
	return nil
}

// Launch Discord App.
func LaunchDiscord(discordVersion DiscordVersion) error {
	var appName string

	switch discordVersion.Index {
	case 0:
		appName = "Discord.app"
	case 1:
		appName = "Discord PTB.app"
	case 2:
		appName = "Discord Canary.app"
	default:
		return fmt.Errorf("unknown Discord version index: %d", discordVersion.Index)
	}

	cmd := exec.Command("open", "-a", appName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch %s: %w", discordVersion.Name, err)
	}

	return nil
}

// Locate the discord_desktop_core module path for the specified Discord version.
func DiscordApplicationSupportPath(discordVersion DiscordVersion) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var discordDir string
	switch discordVersion.Index {
	case 0:
		discordDir = "discord"
	case 1:
		discordDir = "discordptb"
	case 2:
		discordDir = "discordcanary"
	default:
		return "", fmt.Errorf("unknown Discord version")
	}

	basePath := filepath.Join(home, "Library", "Application Support", discordDir)
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return "", fmt.Errorf("discord not found or not installed: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), ".") {
			versions = append(versions, entry.Name())
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no Discord version directories found")
	}

	sort.Strings(versions)
	latestVersion := versions[len(versions)-1]
	modulesPath := filepath.Join(basePath, latestVersion, "modules", "discord_desktop_core")
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		return "", fmt.Errorf("discord_desktop_core module not found")
	}
	return modulesPath, nil
}
