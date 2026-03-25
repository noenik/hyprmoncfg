package apply

import (
	"strings"
	"testing"

	"github.com/crmne/hyprmoncfg/internal/hypr"
	"github.com/crmne/hyprmoncfg/internal/profile"
)

func TestCommandsForProfile(t *testing.T) {
	mon := hypr.Monitor{Name: "DP-1", Make: "Dell", Model: "U2720Q", Serial: "A1"}
	p := profile.New("desk", []profile.OutputConfig{{
		Key:       mon.HardwareKey(),
		Name:      "DP-1",
		Enabled:   true,
		Width:     2560,
		Height:    1440,
		Refresh:   144,
		X:         0,
		Y:         0,
		Scale:     1,
		Transform: 0,
	}})

	cmds, err := CommandsForProfile(p, []hypr.Monitor{mon})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmds) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cmds))
	}
	if !strings.Contains(cmds[0], "DP-1,2560x1440@144") {
		t.Fatalf("unexpected command: %s", cmds[0])
	}
}

func TestCommandsForProfileDisable(t *testing.T) {
	mon := hypr.Monitor{Name: "HDMI-A-1", Make: "LG", Model: "27GP850", Serial: "B2"}
	p := profile.New("dock", []profile.OutputConfig{{
		Key:     mon.HardwareKey(),
		Name:    mon.Name,
		Enabled: false,
	}})

	cmds, err := CommandsForProfile(p, []hypr.Monitor{mon})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmds) != 1 || cmds[0] != "HDMI-A-1,disable" {
		t.Fatalf("unexpected disable command: %+v", cmds)
	}
}
