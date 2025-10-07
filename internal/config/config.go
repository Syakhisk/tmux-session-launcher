package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Directories []DirectoryConfig `yaml:"directories"`
}

// DirectoryConfig represents a directory configuration entry
type DirectoryConfig struct {
	Path  string `yaml:"path"`
	Depth int    `yaml:"depth,omitempty"`
}

// defaultConfig returns the default configuration
func defaultConfig() *Config {
	return &Config{
		Directories: []DirectoryConfig{
			{Path: "$HOME/.config", Depth: 1},
			{Path: "$HOME/Documents", Depth: 1},
			{Path: "$HOME/Desktop"},
		},
	}
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	// Use XDG_CONFIG_HOME if set, otherwise fall back to ~/.config
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to $HOME environment variable
			homeDir = os.ExpandEnv("$HOME")
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "tmux-session-launcher.yaml")
}

// Load loads the configuration from the YAML file or creates a default one
func Load() (*Config, error) {
	configPath := getConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		config := defaultConfig()
		if err := Save(config); err != nil {
			// If we can't save, just return the default config
			return config, nil
		}
		return config, nil
	}

	// Read existing config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration to the YAML file
func Save(config *Config) error {
	configPath := getConfigPath()

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfigPath returns the path to the configuration file (for external use)
func GetConfigPath() string {
	return getConfigPath()
}

// AddDirectory adds a directory to the configuration
func AddDirectory(path string, depth int) error {
	config, err := Load()
	if err != nil {
		return err
	}

	// Check if directory already exists
	for _, dir := range config.Directories {
		if dir.Path == path {
			return fmt.Errorf("directory %s already exists in configuration", path)
		}
	}

	// Add new directory
	config.Directories = append(config.Directories, DirectoryConfig{
		Path:  path,
		Depth: depth,
	})

	return Save(config)
}

// RemoveDirectory removes a directory from the configuration
func RemoveDirectory(path string) error {
	config, err := Load()
	if err != nil {
		return err
	}

	// Find and remove directory
	for i, dir := range config.Directories {
		if dir.Path == path {
			config.Directories = append(config.Directories[:i], config.Directories[i+1:]...)
			return Save(config)
		}
	}

	return fmt.Errorf("directory %s not found in configuration", path)
}

// ListDirectories returns all configured directories
func ListDirectories() ([]DirectoryConfig, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}
	return config.Directories, nil
}
