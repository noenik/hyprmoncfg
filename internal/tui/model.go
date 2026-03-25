package tui

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/crmne/hyprmoncfg/internal/apply"
	"github.com/crmne/hyprmoncfg/internal/hypr"
	"github.com/crmne/hyprmoncfg/internal/profile"
)

type uiMode int

const (
	modeMain uiMode = iota
	modeEdit
	modeSave
	modeConfirm
)

type saveSource int

const (
	saveSourceCurrent saveSource = iota
	saveSourceEdit
)

type refreshMsg struct {
	monitors []hypr.Monitor
	profiles []profile.Profile
	err      error
}

type saveMsg struct {
	name string
	err  error
}

type deleteMsg struct {
	name string
	err  error
}

type applyMsg struct {
	target     string
	snapshot   []string
	returnMode uiMode
	err        error
}

type revertMsg struct {
	err    error
	reason string
}

type tickMsg time.Time

type pendingApply struct {
	snapshot   []string
	target     string
	deadline   time.Time
	returnMode uiMode
}

type Model struct {
	client *hypr.Client
	store  *profile.Store
	engine apply.Engine

	styles styles

	mode      uiMode
	focusPane int

	monitors []hypr.Monitor
	profiles []profile.Profile

	selectedMonitor int
	selectedProfile int

	editOutputs []profile.OutputConfig
	editOutput  int
	editField   int

	saveName   string
	saveSource saveSource

	pending *pendingApply

	status    string
	statusErr bool

	width  int
	height int
}

func NewModel(client *hypr.Client, store *profile.Store) Model {
	return Model{
		client: client,
		store:  store,
		engine: apply.Engine{Client: client},
		styles: newStyles(),
		mode:   modeMain,
		status: "Loading monitors and profiles...",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.refreshCmd(), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case refreshMsg:
		if msg.err != nil {
			m.setStatusErr(msg.err.Error())
			return m, nil
		}
		m.monitors = msg.monitors
		m.profiles = msg.profiles
		m.selectedMonitor = clampIndex(m.selectedMonitor, len(m.monitors))
		m.selectedProfile = clampIndex(m.selectedProfile, len(m.profiles))
		if m.mode == modeEdit && len(m.editOutputs) == 0 {
			m.editOutputs = outputsFromMonitors(m.monitors)
		}
		m.setStatusOK(fmt.Sprintf("Loaded %d monitors and %d profiles", len(m.monitors), len(m.profiles)))
		return m, nil

	case saveMsg:
		if msg.err != nil {
			m.setStatusErr(msg.err.Error())
			m.mode = modeMain
			return m, nil
		}
		m.mode = modeMain
		m.saveName = ""
		m.setStatusOK(fmt.Sprintf("Saved profile %q", msg.name))
		return m, m.refreshCmd()

	case deleteMsg:
		if msg.err != nil {
			m.setStatusErr(msg.err.Error())
			return m, nil
		}
		m.setStatusOK(fmt.Sprintf("Deleted profile %q", msg.name))
		return m, m.refreshCmd()

	case applyMsg:
		if msg.err != nil {
			m.setStatusErr(msg.err.Error())
			m.mode = modeMain
			return m, nil
		}
		m.pending = &pendingApply{
			snapshot:   msg.snapshot,
			target:     msg.target,
			deadline:   time.Now().Add(10 * time.Second),
			returnMode: msg.returnMode,
		}
		m.mode = modeConfirm
		m.statusErr = false
		m.status = fmt.Sprintf("Applied %q. Confirm within 10 seconds or it will revert.", msg.target)
		return m, tea.Batch(m.refreshCmd(), tickCmd())

	case revertMsg:
		m.mode = modeMain
		m.pending = nil
		if msg.err != nil {
			m.setStatusErr(fmt.Sprintf("Revert failed: %v", msg.err))
			return m, nil
		}
		m.setStatusOK("Configuration reverted: " + msg.reason)
		return m, m.refreshCmd()

	case tickMsg:
		if m.mode == modeConfirm && m.pending != nil {
			if time.Now().After(m.pending.deadline) {
				snapshot := append([]string(nil), m.pending.snapshot...)
				return m, m.revertCmd(snapshot, "timeout")
			}
			return m, tickCmd()
		}
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeMain:
			return m.updateMainKeys(msg)
		case modeEdit:
			return m.updateEditKeys(msg)
		case modeSave:
			return m.updateSaveKeys(msg)
		case modeConfirm:
			return m.updateConfirmKeys(msg)
		}
	}
	return m, nil
}

func (m Model) updateMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab", "left", "right":
		m.focusPane = 1 - m.focusPane
		return m, nil
	case "up", "k":
		if m.focusPane == 0 {
			m.selectedMonitor = clampIndex(m.selectedMonitor-1, len(m.monitors))
		} else {
			m.selectedProfile = clampIndex(m.selectedProfile-1, len(m.profiles))
		}
		return m, nil
	case "down", "j":
		if m.focusPane == 0 {
			m.selectedMonitor = clampIndex(m.selectedMonitor+1, len(m.monitors))
		} else {
			m.selectedProfile = clampIndex(m.selectedProfile+1, len(m.profiles))
		}
		return m, nil
	case "r":
		return m, m.refreshCmd()
	case "e":
		m.mode = modeEdit
		m.editOutputs = outputsFromMonitors(m.monitors)
		m.editOutput = 0
		m.editField = 0
		if len(m.editOutputs) == 0 {
			m.setStatusErr("No monitors available for editing")
			m.mode = modeMain
		}
		return m, nil
	case "s":
		m.mode = modeSave
		m.saveSource = saveSourceCurrent
		m.saveName = defaultProfileName()
		return m, nil
	case "a", "enter":
		if len(m.profiles) == 0 {
			m.setStatusErr("No profiles available")
			return m, nil
		}
		p := m.profiles[m.selectedProfile]
		snapshot := apply.SnapshotCommands(m.monitors)
		return m, m.applyCmd(p, snapshot, modeMain)
	case "d":
		if len(m.profiles) == 0 {
			m.setStatusErr("No profiles to delete")
			return m, nil
		}
		name := m.profiles[m.selectedProfile].Name
		return m, m.deleteCmd(name)
	}
	return m, nil
}

func (m Model) updateEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = modeMain
		m.editOutputs = nil
		return m, nil
	case "up", "k":
		m.editOutput = clampIndex(m.editOutput-1, len(m.editOutputs))
		return m, nil
	case "down", "j":
		m.editOutput = clampIndex(m.editOutput+1, len(m.editOutputs))
		return m, nil
	case "left", "h":
		m.editField = clampIndex(m.editField-1, len(editFieldNames))
		return m, nil
	case "right", "l", "tab":
		m.editField = clampIndex(m.editField+1, len(editFieldNames))
		return m, nil
	case " ":
		m.toggleEditOutput()
		return m, nil
	case "+", "=":
		m.adjustEditOutput(1)
		return m, nil
	case "-", "_":
		m.adjustEditOutput(-1)
		return m, nil
	case "r":
		m.editOutputs = outputsFromMonitors(m.monitors)
		m.editOutput = clampIndex(m.editOutput, len(m.editOutputs))
		m.editField = clampIndex(m.editField, len(editFieldNames))
		m.setStatusOK("Editor reset to current monitor state")
		return m, nil
	case "s":
		m.mode = modeSave
		m.saveSource = saveSourceEdit
		m.saveName = defaultProfileName()
		return m, nil
	case "a":
		p := profile.New("temporary", append([]profile.OutputConfig(nil), m.editOutputs...))
		snapshot := apply.SnapshotCommands(m.monitors)
		return m, m.applyCmd(p, snapshot, modeEdit)
	}
	return m, nil
}

func (m Model) updateSaveKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = modeMain
		m.saveName = ""
		return m, nil
	case "enter":
		name := strings.TrimSpace(m.saveName)
		if name == "" {
			m.setStatusErr("Profile name cannot be empty")
			return m, nil
		}
		var p profile.Profile
		if m.saveSource == saveSourceEdit {
			p = profile.New(name, append([]profile.OutputConfig(nil), m.editOutputs...))
		} else {
			p = profile.FromMonitors(name, m.monitors)
		}
		return m, m.saveCmd(p)
	case "backspace":
		if len(m.saveName) > 0 {
			m.saveName = m.saveName[:len(m.saveName)-1]
		}
		return m, nil
	default:
		for _, r := range msg.Runes {
			if unicode.IsPrint(r) {
				if len(m.saveName) < 64 {
					m.saveName += string(r)
				}
			}
		}
		return m, nil
	}
}

func (m Model) updateConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.pending == nil {
		m.mode = modeMain
		return m, nil
	}
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "y", "enter":
		m.mode = m.pending.returnMode
		m.setStatusOK(fmt.Sprintf("Kept configuration for %q", m.pending.target))
		m.pending = nil
		return m, m.refreshCmd()
	case "n", "esc":
		snapshot := append([]string(nil), m.pending.snapshot...)
		return m, m.revertCmd(snapshot, "user request")
	}
	return m, nil
}

func (m Model) View() string {
	switch m.mode {
	case modeEdit:
		return m.renderEdit()
	case modeSave:
		return m.renderSavePrompt()
	case modeConfirm:
		return m.renderConfirm()
	default:
		return m.renderMain()
	}
}

func (m Model) renderMain() string {
	title := m.styles.title.Render("hyprmoncfg")
	panes := m.renderPanes()
	status := m.renderStatus()
	help := m.styles.help.Render("tab switch pane • e edit layout • s save current • a/enter apply profile • d delete profile • r refresh • q quit")
	return strings.Join([]string{title, "", panes, "", status, help}, "\n")
}

func (m Model) renderEdit() string {
	title := m.styles.title.Render("hyprmoncfg • Edit Layout")
	if len(m.editOutputs) == 0 {
		return strings.Join([]string{title, "", m.styles.warning.Render("No editable outputs"), m.styles.help.Render("esc back")}, "\n")
	}

	lines := make([]string, 0, len(m.editOutputs)+2)
	lines = append(lines, m.styles.header.Render("Output [enabled width height refresh x y scale transform]"))
	for i, out := range m.editOutputs {
		fieldValues := []string{
			boolText(out.Enabled),
			fmt.Sprintf("%d", out.Width),
			fmt.Sprintf("%d", out.Height),
			fmt.Sprintf("%.2f", out.Refresh),
			fmt.Sprintf("%d", out.X),
			fmt.Sprintf("%d", out.Y),
			fmt.Sprintf("%.2f", out.Scale),
			fmt.Sprintf("%d", out.Transform),
		}
		for fi := range fieldValues {
			if i == m.editOutput && fi == m.editField {
				fieldValues[fi] = m.styles.focused.Render(fieldValues[fi])
			}
		}
		prefix := "  "
		if i == m.editOutput {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%-10s (%s)  [%s]", prefix, out.Name, out.Key, strings.Join(fieldValues, " ")))
	}

	body := strings.Join(lines, "\n")
	status := m.renderStatus()
	help := m.styles.help.Render("up/down output • left/right field • +/- change • space toggle enabled • a apply • s save profile • r reset • esc back")
	return strings.Join([]string{title, "", body, "", status, help}, "\n")
}

func (m Model) renderSavePrompt() string {
	title := m.styles.title.Render("hyprmoncfg • Save Profile")
	source := "current monitor state"
	if m.saveSource == saveSourceEdit {
		source = "edited layout"
	}
	prompt := fmt.Sprintf("Name (%s): %s", source, m.saveName)
	status := m.renderStatus()
	help := m.styles.help.Render("type name • enter save • backspace delete • esc cancel")
	return strings.Join([]string{title, "", prompt, "", status, help}, "\n")
}

func (m Model) renderConfirm() string {
	title := m.styles.title.Render("hyprmoncfg • Confirm Apply")
	if m.pending == nil {
		return title
	}
	remaining := int(time.Until(m.pending.deadline).Seconds())
	if remaining < 0 {
		remaining = 0
	}
	body := m.styles.warning.Render(fmt.Sprintf("Applied %q. Confirm in %ds or revert.", m.pending.target, remaining))
	status := m.renderStatus()
	help := m.styles.help.Render("y keep configuration • n revert")
	return strings.Join([]string{title, "", body, "", status, help}, "\n")
}

func (m Model) renderPanes() string {
	paneWidth := 50
	if m.width > 120 {
		paneWidth = (m.width - 6) / 2
	}
	if paneWidth < 40 {
		paneWidth = 40
	}

	left := m.renderMonitorPane(paneWidth)
	right := m.renderProfilePane(paneWidth)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
}

func (m Model) renderMonitorPane(width int) string {
	lines := make([]string, 0, len(m.monitors)+1)
	lines = append(lines, m.styles.header.Render("Monitors"))
	if len(m.monitors) == 0 {
		lines = append(lines, "(none)")
	} else {
		for i, mon := range m.monitors {
			prefix := "  "
			if i == m.selectedMonitor && m.focusPane == 0 {
				prefix = "> "
			}
			state := "on"
			if mon.Disabled {
				state = "off"
			}
			line := fmt.Sprintf("%s%-8s %s %dx%d@%.1f pos:%dx%d scale:%.2f %s", prefix, mon.Name, state, mon.Width, mon.Height, mon.RefreshRate, mon.X, mon.Y, mon.Scale, mon.HardwareKey())
			lines = append(lines, line)
		}
	}
	body := strings.Join(lines, "\n")
	panel := m.styles.inactive
	if m.focusPane == 0 {
		panel = m.styles.activePane
	}
	return panel.Width(width).Render(body)
}

func (m Model) renderProfilePane(width int) string {
	lines := make([]string, 0, len(m.profiles)+1)
	lines = append(lines, m.styles.header.Render("Profiles"))
	if len(m.profiles) == 0 {
		lines = append(lines, "(none)")
	} else {
		for i, p := range m.profiles {
			prefix := "  "
			if i == m.selectedProfile && m.focusPane == 1 {
				prefix = "> "
			}
			line := fmt.Sprintf("%s%-20s outputs:%d updated:%s", prefix, p.Name, len(p.Outputs), p.UpdatedAt.Local().Format("2006-01-02 15:04"))
			lines = append(lines, line)
		}
	}
	body := strings.Join(lines, "\n")
	panel := m.styles.inactive
	if m.focusPane == 1 {
		panel = m.styles.activePane
	}
	return panel.Width(width).Render(body)
}

func (m Model) renderStatus() string {
	if m.status == "" {
		return ""
	}
	if m.statusErr {
		return m.styles.statusError.Render(m.status)
	}
	return m.styles.statusOK.Render(m.status)
}

var editFieldNames = []string{"enabled", "width", "height", "refresh", "x", "y", "scale", "transform"}

func (m *Model) toggleEditOutput() {
	if len(m.editOutputs) == 0 {
		return
	}
	m.editOutputs[m.editOutput].Enabled = !m.editOutputs[m.editOutput].Enabled
}

func (m *Model) adjustEditOutput(delta int) {
	if len(m.editOutputs) == 0 {
		return
	}
	out := &m.editOutputs[m.editOutput]
	switch m.editField {
	case 0:
		out.Enabled = !out.Enabled
	case 1:
		out.Width = max(320, out.Width+delta*10)
	case 2:
		out.Height = max(200, out.Height+delta*10)
	case 3:
		out.Refresh = clampFloat(out.Refresh+float64(delta)*0.5, 10, 360)
	case 4:
		out.X += delta * 10
	case 5:
		out.Y += delta * 10
	case 6:
		out.Scale = clampFloat(out.Scale+float64(delta)*0.05, 0.25, 4)
	case 7:
		out.Transform += delta
		if out.Transform < 0 {
			out.Transform = 7
		}
		if out.Transform > 7 {
			out.Transform = 0
		}
	}
}

func (m *Model) setStatusErr(msg string) {
	m.status = msg
	m.statusErr = true
}

func (m *Model) setStatusOK(msg string) {
	m.status = msg
	m.statusErr = false
}

func outputsFromMonitors(monitors []hypr.Monitor) []profile.OutputConfig {
	outs := make([]profile.OutputConfig, 0, len(monitors))
	for _, mon := range monitors {
		outs = append(outs, profile.OutputConfig{
			Key:       mon.HardwareKey(),
			Name:      mon.Name,
			Make:      mon.Make,
			Model:     mon.Model,
			Serial:    mon.Serial,
			Enabled:   !mon.Disabled,
			Width:     max(320, mon.Width),
			Height:    max(200, mon.Height),
			Refresh:   clampFloat(mon.RefreshRate, 10, 360),
			X:         mon.X,
			Y:         mon.Y,
			Scale:     clampFloat(mon.Scale, 0.25, 4),
			Transform: mon.Transform,
		})
	}
	return outs
}

func (m Model) refreshCmd() tea.Cmd {
	client := m.client
	store := m.store
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		monitors, err := client.Monitors(ctx)
		if err != nil {
			return refreshMsg{err: err}
		}
		profiles, err := store.List()
		if err != nil {
			return refreshMsg{err: err}
		}
		return refreshMsg{monitors: monitors, profiles: profiles}
	}
}

func (m Model) saveCmd(p profile.Profile) tea.Cmd {
	store := m.store
	return func() tea.Msg {
		if err := store.Save(p); err != nil {
			return saveMsg{err: err}
		}
		return saveMsg{name: p.Name}
	}
}

func (m Model) deleteCmd(name string) tea.Cmd {
	store := m.store
	return func() tea.Msg {
		if err := store.Delete(name); err != nil {
			return deleteMsg{name: name, err: err}
		}
		return deleteMsg{name: name}
	}
}

func (m Model) applyCmd(p profile.Profile, snapshot []string, returnMode uiMode) tea.Cmd {
	client := m.client
	engine := m.engine
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		monitors, err := client.Monitors(ctx)
		if err != nil {
			return applyMsg{target: p.Name, err: err}
		}
		if _, err := engine.Apply(ctx, p, monitors); err != nil {
			return applyMsg{target: p.Name, err: err}
		}
		return applyMsg{target: p.Name, snapshot: snapshot, returnMode: returnMode}
	}
}

func (m Model) revertCmd(snapshot []string, reason string) tea.Cmd {
	engine := m.engine
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		err := engine.Revert(ctx, snapshot)
		return revertMsg{err: err, reason: reason}
	}
}

func defaultProfileName() string {
	return "profile-" + time.Now().Format("20060102-150405")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func boolText(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

func clampIndex(idx, length int) int {
	if length <= 0 {
		return 0
	}
	if idx < 0 {
		return length - 1
	}
	if idx >= length {
		return 0
	}
	return idx
}

func clampFloat(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
