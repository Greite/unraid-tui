package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/config"
	"github.com/Greite/unraid-tui/internal/tui"
	"github.com/Greite/unraid-tui/internal/tui/onboarding"
)

// Set by GoReleaser via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "unraid-tui",
	Short: "Terminal UI for Unraid server management",
	Long:  "A TUI application to monitor and manage your Unraid server from the terminal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.Exists() {
			if err := runOnboarding(); err != nil {
				return err
			}
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}

		client := api.NewClient(cfg.ServerURL, cfg.APIKey)
		m := tui.NewModel(client)
		p := tea.NewProgram(m)

		_, err = p.Run()
		return err
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("unraid-tui %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runOnboarding() error {
	m := onboarding.New()
	p := tea.NewProgram(m)

	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("onboarding error: %w", err)
	}

	final := result.(onboarding.Model)
	if final.Quitting() {
		return fmt.Errorf("configuration annulee")
	}
	if !final.Completed() {
		return fmt.Errorf("configuration incomplete")
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
