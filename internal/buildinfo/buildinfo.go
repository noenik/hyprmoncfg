package buildinfo

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func Summary(name string) string {
	return fmt.Sprintf("%s %s (%s, %s)", name, Version, Commit, Date)
}
