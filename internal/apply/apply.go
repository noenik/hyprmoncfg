package apply

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/crmne/hyprmoncfg/internal/hypr"
	"github.com/crmne/hyprmoncfg/internal/profile"
)

type Engine struct {
	Client *hypr.Client
}

func CommandsForProfile(p profile.Profile, monitors []hypr.Monitor) ([]string, error) {
	if len(monitors) == 0 {
		return nil, fmt.Errorf("no monitors detected")
	}

	byKey := make(map[string]hypr.Monitor, len(monitors))
	byName := make(map[string]hypr.Monitor, len(monitors))
	for _, m := range monitors {
		byKey[m.HardwareKey()] = m
		byName[m.Name] = m
	}

	commands := make([]string, 0, len(p.Outputs))
	for _, out := range p.Outputs {
		mon, ok := byKey[out.Key]
		if !ok && out.Name != "" {
			mon, ok = byName[out.Name]
		}
		if !ok {
			continue
		}
		commands = append(commands, commandForOutput(mon.Name, out))
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("profile %q does not match any connected monitor", p.Name)
	}
	return commands, nil
}

func SnapshotCommands(monitors []hypr.Monitor) []string {
	commands := make([]string, 0, len(monitors))
	for _, m := range monitors {
		if m.Disabled {
			commands = append(commands, fmt.Sprintf("%s,disable", m.Name))
			continue
		}
		out := profile.OutputConfig{
			Enabled:   true,
			Width:     m.Width,
			Height:    m.Height,
			Refresh:   m.RefreshRate,
			X:         m.X,
			Y:         m.Y,
			Scale:     m.Scale,
			Transform: m.Transform,
		}
		commands = append(commands, commandForOutput(m.Name, out))
	}
	return commands
}

func (e Engine) Apply(ctx context.Context, p profile.Profile, monitors []hypr.Monitor) ([]string, error) {
	if e.Client == nil {
		return nil, fmt.Errorf("nil hypr client")
	}
	commands, err := CommandsForProfile(p, monitors)
	if err != nil {
		return nil, err
	}
	if err := e.Client.BatchKeywordMonitor(ctx, commands); err != nil {
		return nil, err
	}
	return commands, nil
}

func (e Engine) Revert(ctx context.Context, commands []string) error {
	if e.Client == nil {
		return fmt.Errorf("nil hypr client")
	}
	if len(commands) == 0 {
		return nil
	}
	return e.Client.BatchKeywordMonitor(ctx, commands)
}

func commandForOutput(name string, out profile.OutputConfig) string {
	if !out.Enabled {
		return fmt.Sprintf("%s,disable", name)
	}

	mode := "preferred"
	if out.Width > 0 && out.Height > 0 {
		if out.Refresh > 0 {
			mode = fmt.Sprintf("%dx%d@%s", out.Width, out.Height, formatFloat(out.Refresh, 3))
		} else {
			mode = fmt.Sprintf("%dx%d", out.Width, out.Height)
		}
	}

	x := out.X
	y := out.Y
	scale := out.Scale
	if scale <= 0 {
		scale = 1.0
	}
	transform := out.Transform
	if transform < 0 || transform > 7 {
		transform = 0
	}

	return fmt.Sprintf("%s,%s,%dx%d,%s,transform,%d", name, mode, x, y, formatFloat(scale, 3), transform)
}

func formatFloat(v float64, precision int) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	s := strconv.FormatFloat(v, 'f', precision, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}
