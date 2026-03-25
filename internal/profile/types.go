package profile

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/crmne/hyprmoncfg/internal/hypr"
)

type OutputConfig struct {
	Key       string  `json:"key"`
	Name      string  `json:"name"`
	Make      string  `json:"make,omitempty"`
	Model     string  `json:"model,omitempty"`
	Serial    string  `json:"serial,omitempty"`
	Enabled   bool    `json:"enabled"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Refresh   float64 `json:"refresh"`
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Scale     float64 `json:"scale"`
	Transform int     `json:"transform"`
}

type Profile struct {
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Outputs   []OutputConfig `json:"outputs"`
}

func New(name string, outputs []OutputConfig) Profile {
	now := time.Now().UTC()
	p := Profile{
		Name:      strings.TrimSpace(name),
		CreatedAt: now,
		UpdatedAt: now,
		Outputs:   append([]OutputConfig(nil), outputs...),
	}
	p.SortOutputs()
	return p
}

func FromMonitors(name string, monitors []hypr.Monitor) Profile {
	outputs := make([]OutputConfig, 0, len(monitors))
	for _, m := range monitors {
		outputs = append(outputs, OutputConfig{
			Key:       m.HardwareKey(),
			Name:      m.Name,
			Make:      m.Make,
			Model:     m.Model,
			Serial:    m.Serial,
			Enabled:   !m.Disabled,
			Width:     m.Width,
			Height:    m.Height,
			Refresh:   m.RefreshRate,
			X:         m.X,
			Y:         m.Y,
			Scale:     m.Scale,
			Transform: m.Transform,
		})
	}
	return New(name, outputs)
}

func (p *Profile) SortOutputs() {
	sort.SliceStable(p.Outputs, func(i, j int) bool {
		if p.Outputs[i].Enabled != p.Outputs[j].Enabled {
			return p.Outputs[i].Enabled
		}
		return p.Outputs[i].Key < p.Outputs[j].Key
	})
}

func (p Profile) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if len(p.Outputs) == 0 {
		return fmt.Errorf("profile has no outputs")
	}
	for i, out := range p.Outputs {
		if out.Key == "" {
			return fmt.Errorf("output %d has empty key", i)
		}
		if out.Enabled {
			if out.Scale <= 0 {
				return fmt.Errorf("output %d has invalid scale", i)
			}
			if out.Width < 0 || out.Height < 0 {
				return fmt.Errorf("output %d has invalid resolution", i)
			}
		}
	}
	return nil
}

func (p Profile) OutputByKey(key string) (OutputConfig, bool) {
	for _, out := range p.Outputs {
		if out.Key == key {
			return out, true
		}
	}
	return OutputConfig{}, false
}

func (p Profile) Keys() []string {
	keys := make([]string, 0, len(p.Outputs))
	for _, out := range p.Outputs {
		keys = append(keys, out.Key)
	}
	sort.Strings(keys)
	return keys
}
