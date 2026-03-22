package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ServerURL string `mapstructure:"server_url"`
	APIKey    string `mapstructure:"api_key"`
}

// Exists returns true if a config file exists and contains both required fields.
// It also checks environment variables.
func Exists() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	v := viper.New()
	v.SetConfigName(".unraid-tui")
	v.SetConfigType("yaml")
	v.AddConfigPath(home)
	v.SetEnvPrefix("UNRAID")
	v.BindEnv("server_url")
	v.BindEnv("api_key")

	_ = v.ReadInConfig()

	return v.GetString("server_url") != "" && v.GetString("api_key") != ""
}

// FilePath returns the path to the config file.
func FilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".unraid-tui.yaml")
}

// Save writes the config to ~/.unraid-tui.yaml with restricted permissions.
func Save(cfg *Config) error {
	path := FilePath()
	content := fmt.Sprintf("server_url: %q\napi_key: %q\n", cfg.ServerURL, cfg.APIKey)
	return os.WriteFile(path, []byte(content), 0600)
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("finding home directory: %w", err)
	}

	viper.SetConfigName(".unraid-tui")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)

	viper.SetEnvPrefix("UNRAID")
	viper.BindEnv("server_url")
	viper.BindEnv("api_key")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required (set in %s or UNRAID_SERVER_URL env var)", filepath.Join(home, ".unraid-tui.yaml"))
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required (set in %s or UNRAID_API_KEY env var)", filepath.Join(home, ".unraid-tui.yaml"))
	}

	return &cfg, nil
}
