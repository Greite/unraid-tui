package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/config"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/tui"
	"github.com/Greite/unraid-tui/internal/tui/onboarding"
)

// Set by GoReleaser via ldflags.
var (
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	serverFlag string
	langFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "unraid-tui",
	Short: "Terminal UI for Unraid server management",
	Long:  "A TUI application to monitor and manage your Unraid server from the terminal.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if langFlag != "" {
			i18n.SetLang(langFlag)
		} else {
			i18n.DetectLang()
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.Exists() {
			if err := runOnboarding(); err != nil {
				return err
			}
		}

		cfg, err := config.LoadServer(serverFlag)
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

var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "List configured servers",
	Run: func(cmd *cobra.Command, args []string) {
		servers := config.ListServers()
		def := config.DefaultServer()
		if len(servers) == 0 {
			fmt.Println("No servers configured. Run unraid-tui to set up.")
			return
		}
		for _, s := range servers {
			marker := "  "
			if s.Name == def {
				marker = "* "
			}
			fmt.Printf("%s%-15s %s\n", marker, s.Name, s.ServerURL)
		}
	},
}

var addServerCmd = &cobra.Command{
	Use:   "add-server",
	Short: "Add a new server configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := onboarding.New()
		p := tea.NewProgram(m)
		result, err := p.Run()
		if err != nil {
			return err
		}
		final := result.(onboarding.Model)
		if final.Quitting() || !final.Completed() {
			return fmt.Errorf("cancelled")
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverFlag, "server", "", "server name to connect to")
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "language (en, fr)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serversCmd)
	rootCmd.AddCommand(addServerCmd)
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
