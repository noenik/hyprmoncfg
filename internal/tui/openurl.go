package tui

import (
	"fmt"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) openURLCmd(label, url string) tea.Cmd {
	opener := m.openURL
	if opener == nil {
		opener = openExternalURL
	}

	return func() tea.Msg {
		return openURLMsg{
			label: label,
			url:   url,
			err:   opener(url),
		}
	}
}

func openExternalURL(url string) error {
	command, args, err := openExternalURLCommand(url)
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", command, url, err)
	}
	return nil
}

func openExternalURLCommand(url string) (string, []string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "open", []string{url}, nil
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", url}, nil
	default:
		return "xdg-open", []string{url}, nil
	}
}
