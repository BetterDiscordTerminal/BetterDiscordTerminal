package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// RunInstallation performs the BetterDiscord installation
func RunUninstall(discordVersion DiscordVersion) error {
	if err := KillDiscord(discordVersion); err != nil {
		return fmt.Errorf("failed to kill Discord: %w", err)
	}

	if err := RemoveShim(); err != nil {
		return fmt.Errorf("failed to remove shim: %w", err)
	}

	if err := LaunchDiscord(discordVersion); err != nil {
		fmt.Printf("Warning: Failed to launch Discord: %v\n", err)
		fmt.Println("Please launch Discord manually.")
	}
	return nil
}

// Changes back the index.js file to its original state.
func RemoveShim() error {
	shimContent := `module.exports = require("./core.asar");`

	corePath, err := DiscordApplicationSupportPath(DiscordVersion{"Discord", "Standard Discord release", 0})
	if err != nil {
		return fmt.Errorf("failed to get Discord Application Support path: %w", err)
	}

	shimPath := filepath.Join(corePath, "index.js")
	if err := os.WriteFile(shimPath, []byte(shimContent), 0644); err != nil {
		return fmt.Errorf("failed to write shim: %w", err)
	}
	return nil
}
