package main

// RunInstallation performs the BetterDiscord installation
func RunRepair(discordVersion DiscordVersion) error {
	RunUninstall(discordVersion)
	RunInstallation(discordVersion)
	return nil
}
