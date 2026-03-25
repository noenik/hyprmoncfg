package config

import (
	"errors"
	"os"
	"path/filepath"
)

const AppName = "hyprmoncfg"

func BaseDir(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, AppName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if home == "" {
		return "", errors.New("unable to resolve home directory")
	}
	return filepath.Join(home, ".config", AppName), nil
}

func EnsureBaseDir(explicit string) (string, error) {
	base, err := BaseDir(explicit)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(base, "profiles"), 0o755); err != nil {
		return "", err
	}
	return base, nil
}

func ProfilesDir(base string) string {
	return filepath.Join(base, "profiles")
}
