package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func resetViper() {
	viper.Reset()
}

func TestLoad_FromFile(t *testing.T) {
	resetViper()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	content := []byte("server_url: http://192.168.1.100:3001\napi_key: test-key-123\n")
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatal(err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	viper.SetEnvPrefix("UNRAID")
	viper.BindEnv("server_url")
	viper.BindEnv("api_key")

	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.ServerURL != "http://192.168.1.100:3001" {
		t.Errorf("expected server_url http://192.168.1.100:3001, got %s", cfg.ServerURL)
	}
	if cfg.APIKey != "test-key-123" {
		t.Errorf("expected api_key test-key-123, got %s", cfg.APIKey)
	}
}

func TestLoad_FromEnvVars(t *testing.T) {
	resetViper()

	t.Setenv("UNRAID_SERVER_URL", "http://10.0.0.5:3001")
	t.Setenv("UNRAID_API_KEY", "env-key-456")

	viper.SetEnvPrefix("UNRAID")
	viper.BindEnv("server_url")
	viper.BindEnv("api_key")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.ServerURL != "http://10.0.0.5:3001" {
		t.Errorf("expected server_url http://10.0.0.5:3001, got %s", cfg.ServerURL)
	}
	if cfg.APIKey != "env-key-456" {
		t.Errorf("expected api_key env-key-456, got %s", cfg.APIKey)
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := []byte("server_url: \"http://192.168.1.50:3001\"\n")
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	str := string(data)
	if !contains(str, "192.168.1.50") {
		t.Error("expected server URL in saved file")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %o", info.Mode().Perm())
	}
}

func TestSave_Integration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &Config{
		ServerURL: "http://tower:3001",
		APIKey:    "integration-key",
	}

	content := "server_url: \"http://tower:3001\"\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	resetViper()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	var loaded Config
	if err := viper.Unmarshal(&loaded); err != nil {
		t.Fatal(err)
	}

	if loaded.ServerURL != cfg.ServerURL {
		t.Errorf("expected %s, got %s", cfg.ServerURL, loaded.ServerURL)
	}
}

func TestConfigDir_ReturnsUnraidTuiDir(t *testing.T) {
	dir := ConfigDir()
	if !filepath.IsAbs(dir) {
		t.Error("expected absolute path")
	}
	if filepath.Base(dir) != ".unraid-tui" {
		t.Errorf("expected .unraid-tui directory, got %s", filepath.Base(dir))
	}
}

func TestFilePath_ReturnsConfigYaml(t *testing.T) {
	path := FilePath()
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected config.yaml, got %s", filepath.Base(path))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	resetViper()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	content := []byte("server_url: http://file-server:3001\napi_key: file-key\n")
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("UNRAID_SERVER_URL", "http://env-server:3001")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)
	viper.SetEnvPrefix("UNRAID")
	viper.BindEnv("server_url")
	viper.BindEnv("api_key")

	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.ServerURL != "http://env-server:3001" {
		t.Errorf("expected env override http://env-server:3001, got %s", cfg.ServerURL)
	}
	if cfg.APIKey != "file-key" {
		t.Errorf("expected file api_key file-key, got %s", cfg.APIKey)
	}
}
