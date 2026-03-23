package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate a shell completion script for unraid-tui.

To load completions:

  bash:
    source <(unraid-tui completion bash)

    # To load completions for each session, execute once:
    # Linux:
    unraid-tui completion bash > /etc/bash_completion.d/unraid-tui
    # macOS:
    unraid-tui completion bash > $(brew --prefix)/etc/bash_completion.d/unraid-tui

  zsh:
    # If shell completion is not already enabled in your environment,
    # you will need to enable it. You can execute the following once:
    echo "autoload -U compinit; compinit" >> ~/.zshrc

    # To load completions for each session, execute once:
    unraid-tui completion zsh > "${fpath[1]}/_unraid-tui"

    # Oh My Zsh:
    unraid-tui completion zsh > ~/.oh-my-zsh/completions/_unraid-tui

    # You will need to start a new shell for this setup to take effect.

  fish:
    unraid-tui completion fish | source

    # To load completions for each session, execute once:
    unraid-tui completion fish > ~/.config/fish/completions/unraid-tui.fish

  powershell:
    unraid-tui completion powershell | Out-String | Invoke-Expression

    # To load completions for each session, execute once:
    unraid-tui completion powershell > unraid-tui.ps1
    # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
