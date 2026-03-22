package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"go.yaml.in/yaml/v3"
)

const (
	configDir      = ".unraid-tui"
	configFileName = "config"
	configFileType = "yaml"
	keyringService = "unraid-tui"
)

type Config struct {
	ServerURL string `mapstructure:"server_url"`
	APIKey    string `mapstructure:"api_key"`
}

type MultiConfig struct {
	Default  string        `yaml:"default"`
	Language string        `yaml:"language,omitempty"`
	Servers  []ServerEntry `yaml:"servers"`
}

type ServerEntry struct {
	Name      string `yaml:"name"`
	ServerURL string `yaml:"server_url"`
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

func keyringUser(serverName string) string {
	if serverName == "" {
		return "api-key"
	}
	return "api-key/" + serverName
}

// Exists returns true if config is complete.
func Exists() bool {
	cfg, err := loadMultiConfig()
	if err != nil || len(cfg.Servers) == 0 {
		// Try legacy single-server
		return existsLegacy()
	}
	s := cfg.getDefault()
	if s == nil || s.ServerURL == "" {
		return false
	}
	if key, err := keyring.Get(keyringService, keyringUser(s.Name)); err == nil && key != "" {
		return true
	}
	return false
}

func existsLegacy() bool {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(ConfigDir())
	v.SetEnvPrefix("UNRAID")
	v.BindEnv("server_url")
	v.BindEnv("api_key")
	_ = v.ReadInConfig()

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
	if key, err := keyring.Get(keyringService, keyringUser("")); err == nil && key != "" {
		return true
	}
	return false
}

// Save stores a server config.
func Save(cfg *Config) error {
	return SaveServer("default", cfg)
}

// SaveServer stores a named server config.
func SaveServer(name string, cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Load or create multi config
	mc, _ := loadMultiConfig()
	if mc == nil {
		mc = &MultiConfig{}
	}

	// Update or add server
	found := false
	for i, s := range mc.Servers {
		if s.Name == name {
			mc.Servers[i].ServerURL = cfg.ServerURL
			found = true
			break
		}
	}
	if !found {
		mc.Servers = append(mc.Servers, ServerEntry{Name: name, ServerURL: cfg.ServerURL})
	}
	if mc.Default == "" || len(mc.Servers) == 1 {
		mc.Default = name
	}

	// Write config file
	data, err := yaml.Marshal(mc)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(FilePath(), data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	// Save API key to keychain
	if err := keyring.Set(keyringService, keyringUser(name), cfg.APIKey); err != nil {
		// Fallback: save in legacy format
		content := fmt.Sprintf("server_url: %q\napi_key: %q\n", cfg.ServerURL, cfg.APIKey)
		os.WriteFile(FilePath(), []byte(content), 0600)
	}

	removeOldConfig()
	return nil
}

// GetLanguage returns the configured language ("" if not set).
func GetLanguage() string {
	mc, err := loadMultiConfig()
	if err != nil {
		return ""
	}
	return mc.Language
}

// SetLanguage persists the language choice.
func SetLanguage(lang string) error {
	mc, _ := loadMultiConfig()
	if mc == nil {
		mc = &MultiConfig{}
	}
	mc.Language = lang
	dir := ConfigDir()
	os.MkdirAll(dir, 0700)
	data, _ := yaml.Marshal(mc)
	return os.WriteFile(FilePath(), data, 0600)
}

// RemoveServer removes a named server from the config.
func RemoveServer(name string) error {
	mc, err := loadMultiConfig()
	if err != nil {
		return err
	}

	newServers := make([]ServerEntry, 0, len(mc.Servers))
	for _, s := range mc.Servers {
		if s.Name != name {
			newServers = append(newServers, s)
		}
	}
	mc.Servers = newServers

	if mc.Default == name && len(mc.Servers) > 0 {
		mc.Default = mc.Servers[0].Name
	}

	data, _ := yaml.Marshal(mc)
	if err := os.WriteFile(FilePath(), data, 0600); err != nil {
		return err
	}

	// Remove API key from keychain
	keyring.Delete(keyringService, keyringUser(name))
	return nil
}

// Load reads config for the default server.
func Load() (*Config, error) {
	return LoadServer("")
}

// LoadServer reads config for a named server ("" = default).
func LoadServer(name string) (*Config, error) {
	// Try multi-server config
	mc, err := loadMultiConfig()
	if err == nil && len(mc.Servers) > 0 {
		var s *ServerEntry
		if name == "" {
			s = mc.getDefault()
		} else {
			s = mc.getServer(name)
		}
		if s != nil && s.ServerURL != "" {
			apiKey, _ := keyring.Get(keyringService, keyringUser(s.Name))
			if apiKey != "" {
				return &Config{ServerURL: s.ServerURL, APIKey: apiKey}, nil
			}
		}
	}

	// Fallback to legacy
	return loadLegacy()
}

// ListServers returns all configured servers.
func ListServers() []ServerEntry {
	mc, err := loadMultiConfig()
	if err != nil || len(mc.Servers) == 0 {
		// Try legacy
		cfg, err := loadLegacy()
		if err != nil {
			return nil
		}
		return []ServerEntry{{Name: "default", ServerURL: cfg.ServerURL}}
	}
	return mc.Servers
}

// DefaultServer returns the name of the default server.
func DefaultServer() string {
	mc, _ := loadMultiConfig()
	if mc != nil && mc.Default != "" {
		return mc.Default
	}
	return "default"
}

// SetDefault sets the default server.
func SetDefault(name string) error {
	mc, err := loadMultiConfig()
	if err != nil {
		return err
	}
	mc.Default = name
	data, _ := yaml.Marshal(mc)
	return os.WriteFile(FilePath(), data, 0600)
}

func loadMultiConfig() (*MultiConfig, error) {
	data, err := os.ReadFile(FilePath())
	if err != nil {
		return nil, err
	}
	var mc MultiConfig
	if err := yaml.Unmarshal(data, &mc); err != nil {
		return nil, err
	}
	return &mc, nil
}

func (mc *MultiConfig) getDefault() *ServerEntry {
	return mc.getServer(mc.Default)
}

func (mc *MultiConfig) getServer(name string) *ServerEntry {
	for i, s := range mc.Servers {
		if s.Name == name {
			return &mc.Servers[i]
		}
	}
	if len(mc.Servers) > 0 {
		return &mc.Servers[0]
	}
	return nil
}

func loadLegacy() (*Config, error) {
	viper.Reset()
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(ConfigDir())
	viper.SetEnvPrefix("UNRAID")
	viper.BindEnv("server_url")
	viper.BindEnv("api_key")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if migrated := migrateOldConfig(); !migrated {
				return nil, fmt.Errorf("config not found")
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
		return nil, fmt.Errorf("server_url is required")
	}

	if cfg.APIKey == "" {
		if key, err := keyring.Get(keyringService, keyringUser("")); err == nil && key != "" {
			cfg.APIKey = key
		}
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	return &cfg, nil
}

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
	dir := ConfigDir()
	os.MkdirAll(dir, 0700)
	if err := os.WriteFile(FilePath(), data, 0600); err != nil {
		return false
	}
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(dir)
	if err := viper.ReadInConfig(); err != nil {
		return false
	}
	os.Remove(oldPath)
	return true
}

func removeOldConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	os.Remove(filepath.Join(home, ".unraid-tui.yaml"))
}
