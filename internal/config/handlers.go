package config

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// HandlerShowConfigPath shows the path to the configuration file
func HandlerShowConfigPath(ctx context.Context, cmd *cli.Command) error {
	fmt.Println(GetConfigPath())
	return nil
}

// HandlerEditConfig opens the configuration file in the default editor
func HandlerEditConfig(ctx context.Context, cmd *cli.Command) error {
	configPath := GetConfigPath()

	// Ensure config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := defaultConfig()
		if err := Save(config); err != nil {
			return fmt.Errorf("failed to create default configuration: %w", err)
		}
		fmt.Printf("Created default configuration at %s\n", configPath)
	}

	// Determine editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // fallback editor
	}

	// Open editor
	execCmd := exec.Command(editor, configPath)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

// HandlerInitConfig creates a default configuration file
func HandlerInitConfig(ctx context.Context, cmd *cli.Command) error {
	configPath := GetConfigPath()

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration file already exists at %s", configPath)
	}

	// Create default config
	config := defaultConfig()
	if err := Save(config); err != nil {
		return fmt.Errorf("failed to create configuration: %w", err)
	}

	fmt.Printf("Created default configuration at %s\n", configPath)
	return nil
}

// HandlerValidateConfig validates the configuration file
func HandlerValidateConfig(ctx context.Context, cmd *cli.Command) error {
	configPath := GetConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found at %s", configPath)
	}

	// Try to load config
	_, err := Load()
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Println("Configuration is valid")
	return nil
}

// HandlerAddDirectory adds a directory to the configuration
func HandlerAddDirectory(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	depth := cmd.Int("depth")

	var path string
	if len(args) == 0 {
		// Use current directory
		var err error
		path, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	} else {
		// Use specified path
		path = args[0]
		// Convert to absolute path
		if !filepath.IsAbs(path) {
			var err error
			path, err = filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to convert to absolute path: %w", err)
			}
		}
	}

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", path)
	}

	if err := AddDirectory(path, depth); err != nil {
		return err
	}

	fmt.Printf("Added directory: %s (depth: %d)\n", path, depth)
	return nil
}

// HandlerRemoveDirectory removes a directory from the configuration
func HandlerRemoveDirectory(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) != 1 {
		return fmt.Errorf("exactly one path argument is required")
	}

	path := args[0]
	// Convert to absolute path for comparison
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to convert to absolute path: %w", err)
		}
	}

	if err := RemoveDirectory(path); err != nil {
		return err
	}

	fmt.Printf("Removed directory: %s\n", path)
	return nil
}

// HandlerListDirectories lists all configured directories
func HandlerListDirectories(ctx context.Context, cmd *cli.Command) error {
	dirs, err := ListDirectories()
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		fmt.Println("No directories configured")
		return nil
	}

	fmt.Println("Configured directories:")
	for _, dir := range dirs {
		if dir.Depth > 0 {
			fmt.Printf("  %s (depth: %d)\n", dir.Path, dir.Depth)
		} else {
			fmt.Printf("  %s\n", dir.Path)
		}
	}
	return nil
}
