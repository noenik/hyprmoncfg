# hyprmoncfg

A terminal-first monitor configurator for Hyprland, built in Go.

It ships two binaries:

- `hyprmoncfg`: interactive TUI + CLI for monitor profiles.
- `hyprmoncfgd`: daemon that auto-applies best matching profile on monitor hotplug events.

Both binaries expose build metadata with `version`, which is populated in release builds.

## Why this project

This aims to cover a Monique-like workflow without Python runtime dependencies, while keeping profile switching robust in long-running sessions.

## Features

- Save and manage named monitor profiles.
- Apply profiles with safe confirm-or-revert guard (10 second timeout).
- Interactive TUI editor for layout fields:
  - enabled/disabled
  - width/height/refresh
  - x/y position
  - scale
  - transform
- Automatic profile switching daemon:
  - listens to Hyprland `socket2` monitor events
  - polling fallback for resilience
  - debounced apply to avoid hotplug flapping
- Profiles keyed by hardware identity (`make|model|serial`) with connector-name fallback.

## Build

```bash
go build -o bin/hyprmoncfg ./cmd/hyprmoncfg
go build -o bin/hyprmoncfgd ./cmd/hyprmoncfgd
```

## Release pipeline

GitHub Actions ships:

- CI in `.github/workflows/ci.yml`
- tagged releases in `.github/workflows/release.yml`

Create a release by pushing a semver tag such as `v0.1.0`. Goreleaser will build Linux `amd64` and `arm64` archives, attach checksums, and include the packaged systemd units.

## CLI usage

```bash
hyprmoncfg                 # open TUI
hyprmoncfg monitors        # list current outputs
hyprmoncfg profiles        # list saved profiles
hyprmoncfg save desk       # save current layout as profile "desk"
hyprmoncfg apply desk      # apply profile with 10s confirm/revert
hyprmoncfg delete desk     # delete a profile
```

### Apply options

```bash
hyprmoncfg apply desk --confirm-timeout 0
```

Set `--confirm-timeout 0` to disable rollback prompt for non-interactive or scripted usage.

## TUI keybindings

### Main screen

- `tab`: switch focus between monitors and profiles
- `up/down` or `j/k`: move selection
- `e`: open layout editor
- `s`: save current monitor state to new profile
- `a` or `enter`: apply selected profile
- `d`: delete selected profile
- `r`: refresh monitor/profile data
- `q`: quit

### Edit screen

- `up/down`: select output
- `left/right` or `tab`: select field
- `+/-`: adjust selected field
- `space`: toggle output enabled
- `a`: apply edited layout (temporary)
- `s`: save edited layout as profile
- `r`: reset editor from current state
- `esc`: back to main screen

### Confirm screen

- `y` or `enter`: keep configuration
- `n` or `esc`: revert configuration

## Daemon usage

```bash
hyprmoncfgd
```

Options:

```bash
hyprmoncfgd --debounce 1500ms --poll-interval 5s
hyprmoncfgd --profile desk     # force one profile instead of auto-matching
hyprmoncfgd --quiet
```

## Configuration and profile files

Default config directory:

- `~/.config/hyprmoncfg`

Profiles are stored as JSON in:

- `~/.config/hyprmoncfg/profiles/*.json`

Override with:

```bash
hyprmoncfg --config-dir /path/to/config ...
hyprmoncfgd --config-dir /path/to/config ...
```

## Systemd user service (optional)

For packaged installs, the user unit is installed to `/usr/lib/systemd/user/hyprmoncfgd.service`.

For manual local installs to `~/.local/bin`, use `packaging/systemd/hyprmoncfgd.local.service`:

```bash
mkdir -p ~/.config/systemd/user
cp packaging/systemd/hyprmoncfgd.local.service ~/.config/systemd/user/hyprmoncfgd.service
systemctl --user daemon-reload
systemctl --user enable --now hyprmoncfgd
```

## Arch and Omarchy packaging

Recipes live under `packaging/arch`:

- `packaging/arch/hyprmoncfg`: stable package from tagged releases
- `packaging/arch/hyprmoncfg-git`: rolling package from Git

Examples:

```bash
cd packaging/arch/hyprmoncfg-git
makepkg -si
```

```bash
cd packaging/arch/hyprmoncfg
makepkg -si
```

For Omarchy, the `hyprmoncfg-git` package is the better fit while the project is still evolving quickly.

## Development

```bash
go test ./...
```

## Notes

- Requires a running Hyprland session and `hyprctl` in `PATH`.
- Apply operations use `hyprctl keyword monitor ...` commands.
