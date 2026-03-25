package hypr

import (
	"fmt"
	"strings"
)

type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Monitor struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Make            string    `json:"make"`
	Model           string    `json:"model"`
	Serial          string    `json:"serial"`
	Width           int       `json:"width"`
	Height          int       `json:"height"`
	RefreshRate     float64   `json:"refreshRate"`
	X               int       `json:"x"`
	Y               int       `json:"y"`
	Scale           float64   `json:"scale"`
	Transform       int       `json:"transform"`
	Focused         bool      `json:"focused"`
	DPMSStatus      bool      `json:"dpmsStatus"`
	Disabled        bool      `json:"disabled"`
	ActiveWorkspace Workspace `json:"activeWorkspace"`
}

func (m Monitor) IsInternal() bool {
	n := strings.ToLower(m.Name)
	return strings.HasPrefix(n, "edp") || strings.HasPrefix(n, "lvds") || strings.HasPrefix(n, "dsi")
}

func (m Monitor) HardwareKey() string {
	parts := []string{
		cleanIDPart(m.Make),
		cleanIDPart(m.Model),
		cleanIDPart(m.Serial),
	}
	nonEmpty := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	if len(nonEmpty) > 0 {
		return strings.Join(nonEmpty, "|")
	}
	if d := cleanIDPart(m.Description); d != "" {
		return d
	}
	return cleanIDPart(m.Name)
}

func (m Monitor) ModeString() string {
	if m.Width == 0 || m.Height == 0 {
		return "preferred"
	}
	if m.RefreshRate <= 0 {
		return fmt.Sprintf("%dx%d", m.Width, m.Height)
	}
	return fmt.Sprintf("%dx%d@%.3f", m.Width, m.Height, m.RefreshRate)
}

func cleanIDPart(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	v = strings.Join(strings.Fields(v), " ")
	return v
}
