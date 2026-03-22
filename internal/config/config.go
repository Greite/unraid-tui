package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

const (
	configDir      = ".unraid-tui"
	configFileName = "config"
	configFileType = "yaml"
	keyringService = "unraid-tui"
	keyringUser    = "api-key"
)

type Config struct {
	ServerURL string `mapstructure:"server_url"`
	APIKey    string `mapstructure:"api_key"`
}

// ConfigDir returns the path to ~/.unraid-tui/
func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDir)
}

// FilePath returns the path to ~/.unraid-tui/config.yaml
func FilePath() string {
	return filepath.Join(ConfigDir(), configFileName+"."+configFileType)
}

// Exists returns true if config is complete (server_url + api_key available).
func Exists() bool {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(ConfigDir())
	v.SetEnvPrefix("UNRAID")
	v.BindEnv("server_url")
	v.BindEnv("api_key")

	_ = v.ReadInConfig()

	// Also check old location for migration
	if v.GetString("server_url") == "" {
		home, _ := os.UserHomeDir()
		v.AddConfigPath(home)
		v.SetConfigName(".unraid-tui")
		_ = v.ReadInConfig()
	}

	if v.GetString("server_url") == "" {
		return false
	}

	if v.GetString("api_key") != "" {
		return true
	}
	if key, err := keyring.Get(keyringService, keyringUser); err == nil && key != "" {
		return true
	}
	return false
}

// Save stores server_url in config.yaml and api_key in system keychain.
func Save(cfg *Config) error {
	// Ensure config directory exists
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	path := FilePath()
	content := fmt.Sprintf("server_url: %q\n", cfg.ServerURL)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	// Save api_key to system keychain
	if err := keyring.Set(keyringService, keyringUser, cfg.APIKey); err != nil {
		// Fallback: save in file if keychain is unavailable
		content = fmt.Sprintf("server_url: %q\napi_key: %q\n", cfg.ServerURL, cfg.APIKey)
		if err2 := os.WriteFile(path, []byte(content), 0600); err2 != nil {
			return fmt.Errorf("writing config file: %w", err2)
		}
	}

	// Clean up old config location
	removeOldConfig()

	return nil
}

// Load reads config from file + keychain + env vars.
func Load() (*Config, error) {
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(ConfigDir())

	viper.SetEnvPrefix("UNRAID")
	viper.BindEnv("server_url")
	viper.BindEnv("api_key")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Try old location for migration
			if migrated := migrateOldConfig(); !migrated {
				return nil, fmt.Errorf("config not found: run unraid-tui to configure or set env vars")
			}
		} else {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required (set in %s or UNRAID_SERVER_URL env var)", FilePath())
	}

	// API key priority: env var > keychain > yaml file
	if cfg.APIKey == "" {
		if key, err := keyring.Get(keyringService, keyringUser); err == nil && key != "" {
			cfg.APIKey = key
		}
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required (set UNRAID_API_KEY env var or run unraid-tui to configure)")
	}

	// Migrate api_key from file to keychain if present
	if viper.GetString("api_key") != "" {
		if err := keyring.Set(keyringService, keyringUser, cfg.APIKey); err == nil {
			// Rewrite file without api_key
			content := fmt.Sprintf("server_url: %q\n", cfg.ServerURL)
			os.WriteFile(FilePath(), []byte(content), 0600)
		}
	}

	return &cfg, nil
}

// migrateOldConfig moves ~/.unraid-tui.yaml to ~/.unraid-tui/config.yaml
func migrateOldConfig() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	oldPath := filepath.Join(home, ".unraid-tui.yaml")
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return false
	}

	// Create new config dir and write
	dir := ConfigDir()
	os.MkdirAll(dir, 0700)
	if err := os.WriteFile(FilePath(), data, 0600); err != nil {
		return false
	}

	// Re-read from new location
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(dir)
	if err := viper.ReadInConfig(); err != nil {
		return false
	}

	// Remove old file
	os.Remove(oldPath)
	return true
}

// removeOldConfig deletes ~/.unraid-tui.yaml if it exists.
func removeOldConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	oldPath := filepath.Join(home, ".unraid-tui.yaml")
	os.Remove(oldPath)
}
