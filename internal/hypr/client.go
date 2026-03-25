package hypr

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type EventType string

const (
	EventMonitorAdded   EventType = "monitoradded"
	EventMonitorRemoved EventType = "monitorremoved"
)

type Event struct {
	Type  EventType
	Value string
	Raw   string
}

type Client struct {
	hyprctl string
}

func NewClient() (*Client, error) {
	path, err := exec.LookPath("hyprctl")
	if err != nil {
		return nil, fmt.Errorf("hyprctl not found in PATH")
	}
	return &Client{hyprctl: path}, nil
}

func (c *Client) Monitors(ctx context.Context) ([]Monitor, error) {
	cmd := exec.CommandContext(ctx, c.hyprctl, "-j", "monitors", "all")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query monitors: %w", err)
	}
	var monitors []Monitor
	if err := json.Unmarshal(out, &monitors); err != nil {
		return nil, fmt.Errorf("failed to decode hyprctl monitors JSON: %w", err)
	}
	return monitors, nil
}

func (c *Client) KeywordMonitor(ctx context.Context, value string) error {
	cmd := exec.CommandContext(ctx, c.hyprctl, "keyword", "monitor", value)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed applying monitor keyword %q: %w (%s)", value, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (c *Client) BatchKeywordMonitor(ctx context.Context, values []string) error {
	if len(values) == 0 {
		return nil
	}
	batch := make([]string, 0, len(values))
	for _, v := range values {
		batch = append(batch, "keyword monitor "+v)
	}
	cmd := exec.CommandContext(ctx, c.hyprctl, "--batch", strings.Join(batch, " ; "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed batch monitor apply: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func Socket2Path() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return "", errors.New("XDG_RUNTIME_DIR is not set")
	}
	sig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if sig == "" {
		return "", errors.New("HYPRLAND_INSTANCE_SIGNATURE is not set")
	}
	return filepath.Join(runtimeDir, "hypr", sig, ".socket2.sock"), nil
}

func (c *Client) SubscribeMonitorEvents(ctx context.Context) (<-chan Event, <-chan error) {
	events := make(chan Event)
	errorsCh := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errorsCh)

		socketPath, err := Socket2Path()
		if err != nil {
			errorsCh <- err
			return
		}

		dialer := net.Dialer{Timeout: 5 * time.Second}
		conn, err := dialer.DialContext(ctx, "unix", socketPath)
		if err != nil {
			errorsCh <- fmt.Errorf("failed to connect to hyprland socket2: %w", err)
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			event, ok := parseEvent(line)
			if !ok {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case events <- event:
			}
		}

		if err := scanner.Err(); err != nil && !errors.Is(err, context.Canceled) {
			errorsCh <- err
		}
	}()

	return events, errorsCh
}

func parseEvent(line string) (Event, bool) {
	parts := strings.SplitN(line, ">>", 2)
	if len(parts) != 2 {
		return Event{}, false
	}
	typeName := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch EventType(typeName) {
	case EventMonitorAdded, EventMonitorRemoved:
		return Event{Type: EventType(typeName), Value: value, Raw: line}, true
	default:
		return Event{}, false
	}
}
