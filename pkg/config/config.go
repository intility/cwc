package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	configFileName        = "cwc.yaml" // The name of the config file we want to save
	configFilePermissions = 0o600      // The permissions we want to set on the config file
	APIVersion            = "2024-02-01"
)

var SupportedProviders = []string{ //nolint:gochecknoglobals
	"azure",
	"openai",
}

// SanitizeInput trims whitespaces and newlines from a string.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

type Config struct {
	Provider        string `yaml:"provider"`
	Endpoint        string `yaml:"endpoint"`
	ModelDeployment string `yaml:"modelDeployment"`
	Model           string `yaml:"model"`
	ExcludeGitDir   bool   `yaml:"excludeGitDir"`
	UseGitignore    bool   `yaml:"useGitignore"`
	APIVersion      string `yaml:"apiVersion"`
	// Keep APIKey unexported to avoid accidental exposure
	apiKey string
}

type NewConfigParams struct {
	Provider        string
	Endpoint        string
	APIVersion      string
	ModelDeployment string
	Model           string
}

// NewConfig creates a new Config object.
func NewConfig(params NewConfigParams) *Config {
	return &Config{
		Provider:        params.Provider,
		Endpoint:        params.Endpoint,
		APIVersion:      params.APIVersion,
		ModelDeployment: params.ModelDeployment,
		Model:           params.Model,
		ExcludeGitDir:   true,
		UseGitignore:    true,
		apiKey:          "",
	}
}

// SetAPIKey sets the confidential field apiKey.
func (c *Config) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// APIKey returns the confidential field apiKey.
func (c *Config) APIKey() string {
	return c.apiKey
}

func GetConfigDir() (string, error) {
	return XdgConfigPath()
}

func DefaultConfigPath() (string, error) {
	cfgPath, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cfgPath, configFileName), nil
}

func IsWSL() bool {
	_, exists := os.LookupEnv("WSL_DISTRO_NAME")
	return exists
}
