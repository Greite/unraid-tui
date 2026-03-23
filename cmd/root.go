package cmd

import (
	"fmt"
	"log/slog"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/Greite/unraid-tui/internal/api"
	"github.com/Greite/unraid-tui/internal/config"
	"github.com/Greite/unraid-tui/internal/i18n"
	"github.com/Greite/unraid-tui/internal/logging"
	"github.com/Greite/unraid-tui/internal/tui"
	"github.com/Greite/unraid-tui/internal/tui/onboarding"
)

// Set by GoReleaser via ldflags.
var (
	version     = "dev"
	commit      = "none"
	date        = "unknown"
	serverFlag  string
	langFlag    string
	versionFlag bool
	closeLog    func()
)

var rootCmd = &cobra.Command{
	Use:     "unraid-tui",
	Version: "dev",
	Short:   "Terminal UI for Unraid server management",
	Long:    "A TUI application to monitor and manage your Unraid server from the terminal.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		closeLog = logging.Init(config.ConfigDir())
		slog.Info("starting unraid-tui", "version", version, "commit", commit)

		if langFlag != "" {
			i18n.SetLang(langFlag)
		} else if saved := config.GetLanguage(); saved != "" {
			i18n.SetLang(saved)
		} else {
			i18n.DetectLang()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		slog.Info("shutting down")
		if closeLog != nil {
			closeLog()
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
			slog.Error("config load failed", "error", err)
			return fmt.Errorf("configuration error: %w", err)
		}

		client := api.NewClient(cfg.ServerURL, cfg.APIKey)
		m := tui.NewModel(client)
		p := tea.NewProgram(m)

		_, err = p.Run()
		return err
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
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	rootCmd.PersistentFlags().StringVar(&serverFlag, "server", "", "server name to connect to")
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "language (en, fr, zh, hi, es, ar)")

	// Dynamic completion for --server flag
	rootCmd.RegisterFlagCompletionFunc("server", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		servers := config.ListServers()
		names := make([]string, 0, len(servers))
		for _, s := range servers {
			names = append(names, s.Name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	})

	// Static completion for --lang flag
	rootCmd.RegisterFlagCompletionFunc("lang", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return i18n.SupportedLanguages, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(serversCmd)
	rootCmd.AddCommand(addServerCmd)
	rootCmd.AddCommand(completionCmd)
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
