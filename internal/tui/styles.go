package tui

import "github.com/charmbracelet/lipgloss"

type palette struct {
	text               string
	titleFg            string
	titleBg            string
	subtitle           string
	header             string
	subtle             string
	label              string
	value              string
	field              string
	fieldSelectedFg    string
	fieldSelectedBg    string
	groupTitle         string
	paneBorder         string
	paneActiveBorder   string
	tabBorder          string
	tabActiveBorder    string
	tabInactiveFg      string
	tabActiveFg        string
	tabActiveBg        string
	statusOK           string
	statusError        string
	help               string
	footerWarm         string
	footerAccent       string
	footerVersion      string
	warning            string
	badgeAccentFg      string
	badgeAccentBg      string
	badgeOnFg          string
	badgeOnBg          string
	badgeOffFg         string
	badgeOffBg         string
	badgeMutedFg       string
	badgeMutedBg       string
	modalBorder        string
	modalBg            string
	modalTitle         string
	selectedDesc       string
	panelBg            string
	canvasBg           string
	canvasGrid         string
	canvasAxis         string
	cardBorder         string
	cardBg             string
	cardFg             string
	cardMuted          string
	cardDisabledBorder string
	cardDisabledBg     string
	cardDisabledFg     string
	cardDisabledMuted  string
	cardSelectedBorder string
	cardSelectedBg     string
	cardSelectedFg     string
	cardSelectedMuted  string
	snapHighlight      string
}

type styles struct {
	palette          palette
	app              lipgloss.Style
	title            lipgloss.Style
	subtitle         lipgloss.Style
	header           lipgloss.Style
	subtle           lipgloss.Style
	label            lipgloss.Style
	value            lipgloss.Style
	field            lipgloss.Style
	fieldSelected    lipgloss.Style
	group            lipgloss.Style
	groupTitle       lipgloss.Style
	focused          lipgloss.Style
	activePane       lipgloss.Style
	inactivePane     lipgloss.Style
	tabActive        lipgloss.Style
	tabInactive      lipgloss.Style
	statusOK         lipgloss.Style
	statusError      lipgloss.Style
	help             lipgloss.Style
	selectedDesc     lipgloss.Style
	footerLinkWarm   lipgloss.Style
	footerLinkAccent lipgloss.Style
	footerVersion    lipgloss.Style
	warning          lipgloss.Style
	badgeAccent      lipgloss.Style
	badgeOn          lipgloss.Style
	badgeOff         lipgloss.Style
	badgeMuted       lipgloss.Style
	modalBackdrop    lipgloss.Style
	modal            lipgloss.Style
	modalTitle       lipgloss.Style
	canvas           lipgloss.Style
}

func newStyles() styles {
	p := newPalette()

	return styles{
		palette:          p,
		app:              withFG(lipgloss.NewStyle().Padding(0, 1), p.text),
		title:            withBG(withFG(lipgloss.NewStyle().Bold(true).Padding(0, 1), p.titleFg), p.titleBg),
		subtitle:         withFG(lipgloss.NewStyle(), p.subtitle),
		header:           withFG(lipgloss.NewStyle().Bold(true), p.header),
		subtle:           withFG(lipgloss.NewStyle(), p.subtle),
		label:            withFG(lipgloss.NewStyle(), p.label),
		value:            withFG(lipgloss.NewStyle(), p.value),
		field:            withFG(lipgloss.NewStyle().Padding(0, 1), p.field),
		fieldSelected:    withBG(withFG(lipgloss.NewStyle().Padding(0, 1).Bold(true), p.fieldSelectedFg), p.fieldSelectedBg),
		group:            withBG(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.paneBorder)).Padding(0, 1), p.panelBg),
		groupTitle:       withFG(lipgloss.NewStyle().Bold(true), p.groupTitle),
		focused:          withBG(withFG(lipgloss.NewStyle().Bold(true).Padding(0, 1), p.fieldSelectedFg), p.fieldSelectedBg),
		activePane:       withBG(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.paneActiveBorder)).Padding(0, 1), p.panelBg),
		inactivePane:     withBG(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.paneBorder)).Padding(0, 1), p.panelBg),
		tabActive:        withBG(withFG(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.tabActiveBorder)).Padding(0, 1).Bold(true), p.tabActiveFg), p.tabActiveBg),
		tabInactive:      withFG(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.tabBorder)).Padding(0, 1), p.tabInactiveFg),
		statusOK:         withFG(lipgloss.NewStyle().Bold(true), p.statusOK),
		statusError:      withFG(lipgloss.NewStyle().Bold(true), p.statusError),
		help:             withFG(lipgloss.NewStyle(), p.help),
		selectedDesc:     withFG(lipgloss.NewStyle(), p.selectedDesc),
		footerLinkWarm:   withFG(lipgloss.NewStyle().Underline(true), p.footerWarm),
		footerLinkAccent: withFG(lipgloss.NewStyle().Underline(true), p.footerAccent),
		footerVersion:    withFG(lipgloss.NewStyle(), p.footerVersion),
		warning:          withFG(lipgloss.NewStyle().Bold(true), p.warning),
		badgeAccent:      withBG(withFG(lipgloss.NewStyle().Padding(0, 1).Bold(true), p.badgeAccentFg), p.badgeAccentBg),
		badgeOn:          withBG(withFG(lipgloss.NewStyle().Padding(0, 1).Bold(true), p.badgeOnFg), p.badgeOnBg),
		badgeOff:         withBG(withFG(lipgloss.NewStyle().Padding(0, 1), p.badgeOffFg), p.badgeOffBg),
		badgeMuted:       withBG(withFG(lipgloss.NewStyle().Padding(0, 1), p.badgeMutedFg), p.badgeMutedBg),
		modalBackdrop:    lipgloss.NewStyle().Padding(0, 1),
		modal:            withBG(lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color(p.modalBorder)).Padding(1, 2), p.modalBg),
		modalTitle:       withFG(lipgloss.NewStyle().Bold(true), p.modalTitle),
		canvas:           withBG(lipgloss.NewStyle().Padding(0), p.canvasBg),
	}
}

func newPalette() palette {
	if lipgloss.HasDarkBackground() {
		return palette{
			text:               "#E5E7EB",
			titleFg:            "#F3E7D7",
			titleBg:            "#3B2B22",
			subtitle:           "#8E96A8",
			header:             "#F4F6FB",
			subtle:             "#8A93A5",
			label:              "#9CA8BA",
			value:              "#F4F6FB",
			field:              "#E5E7EB",
			fieldSelectedFg:    "#F8FBFF",
			fieldSelectedBg:    "#365E8B",
			groupTitle:         "#D2A56E",
			paneBorder:         "#374151",
			paneActiveBorder:   "#5E86BF",
			tabBorder:          "#374151",
			tabActiveBorder:    "#5E86BF",
			tabInactiveFg:      "#A0A8B8",
			tabActiveFg:        "#F8FBFF",
			tabActiveBg:        "#243C5A",
			statusOK:           "#8BCF63",
			statusError:        "#E67C73",
			help:               "#768093",
			footerWarm:         "#D6A948",
			footerAccent:       "#C792EA",
			footerVersion:      "#8E96A8",
			warning:            "#D6A948",
			badgeAccentFg:      "#111827",
			badgeAccentBg:      "#D2A56E",
			badgeOnFg:          "#0F1A10",
			badgeOnBg:          "#8BCF63",
			badgeOffFg:         "#D7DCE5",
			badgeOffBg:         "#475266",
			badgeMutedFg:       "#E5E7EB",
			badgeMutedBg:       "#64748B",
			modalBorder:        "#7EA8FF",
			modalBg:            "#131924",
			modalTitle:         "#F3E7D7",
			selectedDesc:       "#DDEAFF",
			panelBg:            "",
			canvasBg:           "",
			canvasGrid:         "#273142",
			canvasAxis:         "#364153",
			cardBorder:         "#6B778C",
			cardBg:             "",
			cardFg:             "#E9EDF5",
			cardMuted:          "#BAC2D0",
			cardDisabledBorder: "#8F6A72",
			cardDisabledBg:     "",
			cardDisabledFg:     "#D7C1C5",
			cardDisabledMuted:  "#B79FA4",
			cardSelectedBorder: "#8EC5FF",
			cardSelectedBg:     "#355C8A",
			cardSelectedFg:     "#F8FBFF",
			cardSelectedMuted:  "#DDEAFF",
			snapHighlight:      "#F5C86A",
		}
	}

	return palette{
		text:               "#1F2937",
		titleFg:            "#6C4320",
		titleBg:            "#EAD5C2",
		subtitle:           "#667085",
		header:             "#111827",
		subtle:             "#667085",
		label:              "#6B7280",
		value:              "#111827",
		field:              "#1F2937",
		fieldSelectedFg:    "#F8FBFF",
		fieldSelectedBg:    "#2F6FBC",
		groupTitle:         "#9A5B13",
		paneBorder:         "#B8C3D4",
		paneActiveBorder:   "#2F6FBC",
		tabBorder:          "#C8D1DE",
		tabActiveBorder:    "#2F6FBC",
		tabInactiveFg:      "#4B5563",
		tabActiveFg:        "#173A5E",
		tabActiveBg:        "#E4EEF9",
		statusOK:           "#2F7A2F",
		statusError:        "#BE3A34",
		help:               "#64748B",
		footerWarm:         "#A16207",
		footerAccent:       "#7C3AED",
		footerVersion:      "#667085",
		warning:            "#A16207",
		badgeAccentFg:      "#111827",
		badgeAccentBg:      "#E0B15C",
		badgeOnFg:          "#FFFFFF",
		badgeOnBg:          "#4E8F3A",
		badgeOffFg:         "#334155",
		badgeOffBg:         "#D5DCE7",
		badgeMutedFg:       "#334155",
		badgeMutedBg:       "#CBD5E1",
		modalBorder:        "#4A7CC2",
		modalBg:            "#F8FAFC",
		modalTitle:         "#6C4320",
		selectedDesc:       "#2B5E98",
		panelBg:            "",
		canvasBg:           "",
		canvasGrid:         "#D6DDE7",
		canvasAxis:         "#C2CBD8",
		cardBorder:         "#8A99AF",
		cardBg:             "",
		cardFg:             "#1F2937",
		cardMuted:          "#617187",
		cardDisabledBorder: "#B88A8A",
		cardDisabledBg:     "",
		cardDisabledFg:     "#7C5A5A",
		cardDisabledMuted:  "#927878",
		cardSelectedBorder: "#2F6FBC",
		cardSelectedBg:     "#DCEAF9",
		cardSelectedFg:     "#123B65",
		cardSelectedMuted:  "#345A84",
		snapHighlight:      "#B7791F",
	}
}

func withFG(style lipgloss.Style, value string) lipgloss.Style {
	if value == "" {
		return style
	}
	return style.Foreground(lipgloss.Color(value))
}

func withBG(style lipgloss.Style, value string) lipgloss.Style {
	if value == "" {
		return style
	}
	return style.Background(lipgloss.Color(value))
}
