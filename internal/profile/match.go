package profile

import (
	"sort"
	"strings"

	"github.com/crmne/hyprmoncfg/internal/hypr"
)

func MatchScore(p Profile, monitors []hypr.Monitor) int {
	if len(monitors) == 0 || len(p.Outputs) == 0 {
		return 0
	}

	connected := make(map[string]struct{}, len(monitors))
	for _, m := range monitors {
		connected[m.HardwareKey()] = struct{}{}
	}

	profileEnabled := make(map[string]struct{})
	for _, o := range p.Outputs {
		if o.Enabled {
			profileEnabled[o.Key] = struct{}{}
		}
	}
	if len(profileEnabled) == 0 {
		return 0
	}

	intersection := 0
	for key := range connected {
		if _, ok := profileEnabled[key]; ok {
			intersection++
		}
	}
	if intersection == 0 {
		return 0
	}

	missingFromCurrent := len(profileEnabled) - intersection
	unknownCurrent := len(connected) - intersection

	// High reward for overlap, moderate penalty for mismatch.
	return intersection*100 - missingFromCurrent*30 - unknownCurrent*20
}

func BestMatch(profiles []Profile, monitors []hypr.Monitor) (Profile, int, bool) {
	type candidate struct {
		profile Profile
		score   int
	}
	candidates := make([]candidate, 0, len(profiles))
	for _, p := range profiles {
		score := MatchScore(p, monitors)
		if score <= 0 {
			continue
		}
		candidates = append(candidates, candidate{profile: p, score: score})
	}
	if len(candidates) == 0 {
		return Profile{}, 0, false
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return strings.ToLower(candidates[i].profile.Name) < strings.ToLower(candidates[j].profile.Name)
	})
	return candidates[0].profile, candidates[0].score, true
}

func MonitorSetHash(monitors []hypr.Monitor) string {
	if len(monitors) == 0 {
		return "none"
	}
	keys := make([]string, 0, len(monitors))
	for _, m := range monitors {
		keys = append(keys, m.HardwareKey())
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}
